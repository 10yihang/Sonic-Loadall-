# Sonic-Loadall-性能测试
针对LoadAll与ConcurrentRead的性能对比测试

## Big.go是[大Json](https://github.com/bytedance/sonic/blob/main/testdata/twitter.json)与中等大小Json的结构

## Small.go是[小Json](https://github.com/bytedance/sonic/blob/main/testdata/small.go)的结构

## 以下是报告

# Sonic 性能基准测试报告

## 测试环境
- **操作系统**: Windows
- **处理器**: 12th Gen Intel(R) Core(TM) i7-12700
- **Go 版本**: go1.24.4
- **Sonic 版本**: v1.13.3

## 测试数据规模
- **Small JSON**: 1KB 
- **Mid JSON**: 14KB 
- **Big JSON**: 632KB 
- **所有大中小的json都是sonic官方文档提到的，且这里我们只进行只读的测试**

## 研究背景与问题发现 🔍

### **问题起源：Len() 方法的异常行为**

在使用 Sonic 库处理 JSON 数据时，我们发现了一个重要问题：**在不调用 `LoadAll()` 或使用`ConcurrentRead` 的情况下，`Len()` 方法返回的结果不正确**。

#### 🧪 **问题复现代码**
```go
data := `{"data":[{"name":"item1"},{"name":"item2"},{"name":"item3"}]}`
jsondata := []byte(data)

dataRoot, err := sonic.Get(jsondata, "data")
lenBefore, _ := dataRoot.Len()
fmt.Println("Before Index(0) Length:", lenBefore)   // 输出：0 ❌
Index0Interface, _ := dataRoot.Index(0).Interface()
fmt.Println("Test Index 0:", Index0Interface)

lenBefore, _ = dataRoot.Len()
fmt.Println("Before LoadAll() Length:", lenBefore)  // 输出: 1 ❌

dataRoot.LoadAll()
lenAfter, _ := dataRoot.Len()
fmt.Println("After LoadAll() Length:", lenAfter)    // 输出: 3 ✅
```

#### 📊 **测试结果分析**
```
Before Index(0) Length: 0      
Test Index 0: map[name:item1]
Before LoadAll() Length: 1      
After LoadAll() Length: 3       
```

#### 🤔 **问题分析**
- **预期结果**: 数组 `[{"name":"item1"},{"name":"item2"},{"name":"item3"}]` 的长度应该是 **3**
- **实际结果**: 在调用 `LoadAll()` 之前，`Len()` 返回 **1**
- **根本原因**: Sonic 的懒加载机制导致只有在进行解析完整解析后才能获取正确的数组长度
<img src="image.png" style="zoom:50%" />

## 详细性能对比

### 1. 小型 JSON (Small JSON ~1KB)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能提升 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 3,462 | 3,216 | 11 | - |
| ConcurrentRead | 612 | 120 | 5 | **5.7x 更快** |

### 2. 中型 JSON (Mid JSON ~10KB)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能提升 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 13,355 | 15,611 | 11 | - |
| ConcurrentRead | 166 | 192 | 4 | **80.6x 更快** |

### 3. 大型 JSON (Big JSON ~630KB)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能提升 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 526,028 | 662,643 | 118 | - |
| ConcurrentRead | 31,804 | 216 | 5 | **16.5x 更快** |

## 特殊场景分析

### 4. 单一字段访问 (Big JSON)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能提升 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 616,260 | 662,922 | 117 | - |
| ConcurrentRead | 67 | 48 | 1 | **9,198x 更快** |

**重要发现**: 对于单一字段访问，ConcurrentRead 的优势极其明显！

### 5. 重复访问同一路径 50 次 (Big JSON)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能提升 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 497,446 | 663,641 | 117 | - |
| ConcurrentRead | 2,592 | 2,400 | 50 | **192x 更快** |

### 6. 复杂访问模式 (Big JSON)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能提升 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 559,518 | 662,833 | 118 | - |
| ConcurrentRead | 148,035 | 6,816 | 126 | **3.8x 更快** |

### 7. 重复访问所有字段 (Big JSON - 420次字段访问)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) | 性能对比 |
|------|-----------------|----------------|-------------------|---------|
| LoadAll | 497,582 | 663,974 | 118 | **4.4x 更快** |
| ConcurrentRead | 2,191,077 | 18,240 | 420 | - |

**重要发现**: 这是 **LoadAll 首次在执行时间上优于 ConcurrentRead** 的场景！(不过貌似实际场景不会有这种情况)

**原因分析**:
- **访问密度极高**: 20轮 × 21个不同字段 = 420次字段访问
- **LoadAll 优势**: 一次性解析 + 420次零开销访问
- **ConcurrentRead 劣势**: 420次解析路径与获取锁释放锁的累积开销

### 8. 内存分配专项测试 (Big JSON - 访问 20 个字段)
| 方式 | 执行时间 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) |
|------|-----------------|----------------|-------------------|
| LoadAll | 540,668 | 662,247 | 117 |
| ConcurrentRead | 1,002 | 960 | 20 |

## 关键结论 📊

### 1. **性能选择取决于访问模式** ⚖️
- **稀疏访问场景**: ConcurrentRead 显著优于 LoadAll
- **密集访问场景**: LoadAll 可能优于 ConcurrentRead
- **内存效率**: ConcurrentRead 在基本所有只读场景下都大幅优于 LoadAll

### 2. **ConcurrentRead 优势场景** ✅
- **单一或少量字段访问** (**绝对优势**) 
- **重复相同路径**
- **复杂但稀疏的访问**

### 3. **LoadAll 优势场景** ⚡ 
- **高密度多字段访问**: 当访问字段数量接近或超过 JSON 总字段数的大部分时
- **示例**: 420次不同字段访问时，LoadAll 比 ConcurrentRead 快 4.4x

### 4. **JSON 大小影响**
- Small JSON: ConcurrentRead 5.7x 提升
- Mid JSON: ConcurrentRead 80.6x 提升
- Big JSON: 取决于访问模式 (稀疏访问 ConcurrentRead 绝对优势，大量密集访问 LoadAll 可能更好)

### 5. **内存使用效率** 💚
- ConcurrentRead 内存分配极少，根据访问字段进行分配，
- LoadAll 需要分配大量内存，所有字段都需加载
- **ConcurrentRead 内存效率要高于LoadAll**

### ⚖️ **选择决策指南**

#### **访问密度评估**
```go
// 计算访问密度 = 访问字段数 / JSON总字段数
访问密度 < 10%  → 强烈推荐 ConcurrentRead
访问密度 10-50% → 推荐 ConcurrentRead
访问密度 50-80% → 需要基准测试来决定
访问密度 > 80%  → 考虑 LoadAll
```

#### **访问次数评估**
```go
单次访问          → ConcurrentRead
少量多次访问      → ConcurrentRead
密集访问          → 需要测试，LoadAll 可能更好
```

---

## 技术原理深度分析 🔍

### **为什么 ConcurrentRead 性能远超 LoadAll？**

通过分析 Sonic 库的核心源码，我们发现了性能差异的根本原因：

#### 核心解析函数 `parseRaw()` 的执行逻辑

```go
func (self *Node) parseRaw(full bool) {
    lock := self.lock()        // 尝试获取锁
    defer self.unlock()        // 释放锁
    
    if !self.isRaw() {
        return  // 如果已经解析过了，直接返回
    }
    
    raw := self.toString()
    parser := NewParserObj(raw)
    var e types.ParsingError
    
    if full {
        // LoadAll 方式：解析整个 JSON 结构
        parser.noLazy = true
        *self, e = parser.Parse()
    } else if lock {  // ConcurrentRead 方式
        // 只解析当前层级，延迟解析子结构
        var n Node
        parser.noLazy = true
        parser.loadOnce = true  // 👈 关键：只加载一次当前层
        n, e = parser.Parse()
        self.assign(n)
    } else {
        // 普通情况：懒加载
        *self, e = parser.Parse()
    }
    
    if e != 0 {
        *self = *newSyntaxError(parser.syntaxError(e))
    }
}
```

#### 🔑 **关键差异分析**

| 方式 | 解析策略 | parser.loadOnce | 内存占用 | 解析深度 |
|------|---------|----------------|---------|---------|
| **LoadAll** | `full = true` | ❌ 未设置 | 📈 高 - 解析整个结构 | 🌳 递归解析所有层级 |
| **ConcurrentRead** | `lock = true` | ✅ `loadOnce = true` | 📉 低 - 按需解析 | 🍃 仅解析当前访问层级 |

#### 🚨 **Len() 方法问题的技术根源**

**为什么普通 Get 方式的 Len() 返回错误结果？**

```go
// 问题场景
data := `{"data":[{"name":"item1"},{"name":"item2"},{"name":"item3"}]}`
dataRoot, _ := sonic.Get(jsondata, "data")
len, _ := dataRoot.Len()  // 返回 1，而不是期望的 3
```

**根本原因分析**：
1. **懒加载机制**: 普通 `Get()` 采用懒加载，只解析到能访问目标路径的最小程度
2. **数组未完全解析**: 对于数组，懒加载只解析第一个元素来确定数据类型
3. **Len() 依赖完整结构**: `Len()` 方法需要完整的数组结构才能返回正确长度

**三种方式的解析行为对比**：

| 解析方式 | 数组解析程度 | Len()结果 | 技术原理 |
|---------|-------------|-----------|---------|
| **普通 Get** | 🔍 仅解析第1个元素 | ❌ 1 (错误) | 懒加载，最小化解析 |
| **LoadAll** | 🌍 解析所有元素 | ✅ 3 (正确) | 强制完整解析 |
| **ConcurrentRead** | 🎯 按需完整解析访问路径 | ✅ 3 (正确) | 智能按需解析 |

**ConcurrentRead 解决方案的优势**：
- ✅ **正确性**: 能正确返回 `Len()` 结果
- ✅ **效率**: 只解析访问路径上的完整结构，避免过度解析
- ✅ **性能**: 比 LoadAll 快几十到几千倍
- ✅ **内存**: 内存使用最优化

#### 📊 **性能优势的技术原因**

1. **懒加载策略**
   - **ConcurrentRead**: 设置 `loadOnce = true`，只解析当前需要的层级
   - **LoadAll**: 递归解析整个 JSON 树结构，包括不需要的部分

2. **内存分配模式**
   - **ConcurrentRead**: 按需分配内存，访问什么解析什么
   - **LoadAll**: 一次性分配整个 JSON 结构的内存

3. **解析深度控制**
   - **ConcurrentRead**: 延迟解析，深度优先且按需进行
   - **LoadAll**: 广度优先，必须完整解析所有节点

4. **CPU 利用效率**
   - **ConcurrentRead**: 减少不必要的解析计算
   - **LoadAll**: 浪费 CPU 资源解析未使用的数据

#### 🎯 **实际影响示例**

以 630KB 的大型 JSON 为例：

```go
// LoadAll 方式 - 解析整个 630KB
{
    "users": [...1000+ 用户对象...],     // 全部解析 ❌
    "metadata": {...},                  // 全部解析 ❌
    "pagination": {...},                // 全部解析 ❌
    "settings": {...}                   // 全部解析 ❌
}

// ConcurrentRead 方式 - 只解析需要的部分
访问 users[0].name 时：
- 只解析 users 数组的第一个元素 ✅
- 只解析该元素的 name 字段 ✅
- 其他 999+ 用户对象保持未解析状态 ✅
```

**结果**: ConcurrentRead 避免了 99% 以上的不必要解析工作！


## 性能可优化点 ⚡

### **1. 锁机制优化：单线程只读场景的进一步提升**

#### 🔒 **当前 ConcurrentRead 的锁开销分析**

从源码可以看到，ConcurrentRead 方式会执行锁操作：

```go
func (self *Node) parseRaw(full bool) {
    lock := self.lock()        // 📌 获取锁
    defer self.unlock()        // 📌 释放锁
    
    // ... 解析逻辑
    if lock {  // ConcurrentRead 进入这个分支
        var n Node
        parser.noLazy = true
        parser.loadOnce = true
        n, e = parser.Parse()
        self.assign(n)
    }
}
```

#### 📊 **锁开销的理论分析**

| 操作类型 | 场景影响 | 优化潜力 |
|---------|---------|---------|
| **锁获取/释放** | 高频访问时累积明显 | 🟡 中等 |

#### 🎯 **特定场景的优化机会**

**理想优化场景**：
-  只读访问 
-  单线程环境 
-  高频调用 
-  性能敏感 

#### 🛠️ **理论实现方案**

**方案1: 条件锁机制**
```go
// 伪代码示例
func GetWithOptionsOptimized(data []byte, opts SearchOptions, path ...string) {
    if opts.ReadOnly && !opts.EnableConcurrency {
        // 使用无锁版本
        return getWithoutLock(data, path...)
    }
    // 使用标准 ConcurrentRead
    return sonic.GetWithOptions(data, opts, path...)
}
```

**方案2: 编译时优化**
```go
// 编译时确定的优化版本
//go:build single_thread
func parseRawOptimized(self *Node, full bool) {
    // 跳过锁机制的版本
    if !self.isRaw() {
        return
    }
    // 直接解析，无锁开销
}
```
## 结论

通过全面的基准测试和深度技术分析，**ConcurrentRead 方式在基本所有只读场景中都要更优一些**：

### 🎯 **解决方案对比**
| 方式 | Len()结果 | 性能 | 内存使用 | 推荐度 |
|------|----------|------|---------|--------|
| 普通 Get | ❌ 不正确 | 🚀 快 | 💚 低 | ⚠️ 有限制 |
| LoadAll | ✅ 正确 | 🐌 慢 | 🔴 高 | ❌ 仅少数场景推荐 |
| ConcurrentRead | ✅ 正确 | 🚀 快 | 💚 低 | ✅ 强烈推荐 |

#### 文档参考
- [Sonic 文档](https://github.com/bytedance/sonic/blob/main/docs/INTRODUCTION_ZH_CN.md)
- [性能基准测试](https://github.com/bytedance/sonic/blob/main/README_ZH_CN.md)
- [测试代码](https://github.com/10yihang/Sonic-Loadall-)
