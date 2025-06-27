package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

var (
	smallJSONData []byte
	midJSONData   []byte
	bigJSONData   []byte
)

// 初始化测试数据
func init() {
	var err error

	smallJSONData, err = os.ReadFile("small.json")
	if err != nil {
		panic("Failed to read small.json: " + err.Error())
	}

	midJSONData, err = os.ReadFile("mid.json")
	if err != nil {
		panic("Failed to read mid.json: " + err.Error())
	}

	bigJSONData, err = os.ReadFile("big.json")
	if err != nil {
		panic("Failed to read big.json: " + err.Error())
	}
}

// ============= Small JSON 基准测试 =============

// 测试小型 JSON - LoadAll 方式
func BenchmarkSmallJSON_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(smallJSONData)
		if err != nil {
			b.Fatal(err)
		}

		// 预加载所有数据
		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 多次访问不同路径
		root.GetByPath("id")
		root.GetByPath("title")
		root.GetByPath("price")
		root.GetByPath("author", "name")
		root.GetByPath("author", "age")
	}
}

// 测试小型 JSON - ConcurrentRead 方式
func BenchmarkSmallJSON_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 多次访问不同路径，每次都使用 ConcurrentRead
		sonic.GetWithOptions(smallJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "id")
		sonic.GetWithOptions(smallJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "title")
		sonic.GetWithOptions(smallJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "price")
		sonic.GetWithOptions(smallJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "author", "name")
		sonic.GetWithOptions(smallJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "author", "age")
	}
}

// ============= Mid JSON 基准测试 =============

// 测试中型 JSON - LoadAll 方式
func BenchmarkMidJSON_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(midJSONData)
		if err != nil {
			b.Fatal(err)
		}

		// 预加载所有数据
		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 多次访问不同路径
		root.GetByPath("statuses", "0", "text")
		root.GetByPath("statuses", "0", "user", "name")
		root.GetByPath("statuses", "0", "user", "followers_count")
		root.GetByPath("statuses", "0", "created_at")
	}
}

// 测试中型 JSON - ConcurrentRead 方式
func BenchmarkMidJSON_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 多次访问不同路径，每次都使用 ConcurrentRead
		sonic.GetWithOptions(midJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "text")
		sonic.GetWithOptions(midJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "user", "name")
		sonic.GetWithOptions(midJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "user", "followers_count")
		sonic.GetWithOptions(midJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "created_at")
	}
}

// ============= Big JSON 基准测试 =============

// 测试大型 JSON - LoadAll 方式
func BenchmarkBigJSON_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(bigJSONData)
		if err != nil {
			b.Fatal(err)
		}

		// 预加载所有数据
		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 多次访问不同路径
		root.GetByPath("statuses", "0", "text")
		root.GetByPath("statuses", "0", "user", "name")
		root.GetByPath("statuses", "0", "user", "followers_count")
		root.GetByPath("search_metadata", "count")
		root.GetByPath("statuses", "5", "user", "screen_name")
	}
}

// 测试大型 JSON - ConcurrentRead 方式
func BenchmarkBigJSON_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 多次访问不同路径，每次都使用 ConcurrentRead
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "text")
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "user", "name")
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "user", "followers_count")
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "search_metadata", "count")
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "5", "user", "screen_name")
	}
}

// ============= 内存分配测试 =============

// 测试内存分配 - LoadAll 方式 (Big JSON)
func BenchmarkBigJSON_LoadAll_Memory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(bigJSONData)
		if err != nil {
			b.Fatal(err)
		}

		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 访问多个路径
		for j := 0; j < 10; j++ {
			root.GetByPath("statuses", "0", "text")
			root.GetByPath("statuses", "0", "user", "name")
		}
	}
}

// 测试内存分配 - ConcurrentRead 方式 (Big JSON)
func BenchmarkBigJSON_ConcurrentRead_Memory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// 访问多个路径
		for j := 0; j < 10; j++ {
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "text")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "name")
		}
	}
}

// ============= 复杂访问模式测试 =============

// 测试复杂访问模式 - LoadAll 方式
func BenchmarkComplexAccess_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(bigJSONData)
		if err != nil {
			b.Fatal(err)
		}

		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}
		// 模拟复杂的访问模式
		// 1. 获取状态数组的长度
		statuses := root.GetByPath("statuses")
		if statuses.Exists() {
			length, _ := statuses.Len()
			// 2. 访问前几个状态的详细信息
			for j := 0; j < min(length, 3); j++ {
				root.GetByPath("statuses", string(rune(j+'0')), "text")
				root.GetByPath("statuses", string(rune(j+'0')), "user", "name")
				root.GetByPath("statuses", string(rune(j+'0')), "user", "followers_count")
				root.GetByPath("statuses", string(rune(j+'0')), "created_at")
			}
		}

		// 3. 获取搜索元数据
		root.GetByPath("search_metadata", "count")
		root.GetByPath("search_metadata", "completed_in")
	}
}

// 测试复杂访问模式 - ConcurrentRead 方式
func BenchmarkComplexAccess_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 模拟复杂的访问模式
		// 1. 获取状态数组的长度
		statuses, _ := sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses")
		if statuses.Exists() {
			length, _ := statuses.Len()
			// 2. 访问前几个状态的详细信息
			for j := 0; j < min(length, 3); j++ {
				sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
					ConcurrentRead: true,
				}, "statuses", string(rune(j+'0')), "text")
				sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
					ConcurrentRead: true,
				}, "statuses", string(rune(j+'0')), "user", "name")
				sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
					ConcurrentRead: true,
				}, "statuses", string(rune(j+'0')), "user", "followers_count")
				sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
					ConcurrentRead: true,
				}, "statuses", string(rune(j+'0')), "created_at")
			}
		}

		// 3. 获取搜索元数据
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "search_metadata", "count")
		sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "search_metadata", "completed_in")
	}
}

// ============= 解析 vs 部分访问测试 =============

// 测试完整解析 - 标准 JSON 库
func BenchmarkBigJSON_StdJsonUnmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var twitter TwitterStruct
		err := json.Unmarshal(bigJSONData, &twitter)
		if err != nil {
			b.Fatal(err)
		}

		// 访问一些字段来模拟实际使用
		if len(twitter.Statuses) > 0 {
			_ = twitter.Statuses[0].Text
			_ = twitter.Statuses[0].User.Name
		}
		_ = twitter.SearchMetadata.Count
	}
}

// 测试 Sonic 完整解析
func BenchmarkBigJSON_SonicUnmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var twitter TwitterStruct
		err := sonic.Unmarshal(bigJSONData, &twitter)
		if err != nil {
			b.Fatal(err)
		}

		// 访问一些字段来模拟实际使用
		if len(twitter.Statuses) > 0 {
			_ = twitter.Statuses[0].Text
			_ = twitter.Statuses[0].User.Name
		}
		_ = twitter.SearchMetadata.Count
	}
}

// ============= 单一访问测试 =============

// 测试单一字段访问 - LoadAll 方式
func BenchmarkSingleAccess_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(bigJSONData)
		if err != nil {
			b.Fatal(err)
		}

		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 只访问一个字段
		node := root.GetByPath("statuses", "0", "text")
		if node.Exists() {
			_, _ = node.String()
		}
	}
}

// 测试单一字段访问 - ConcurrentRead 方式
func BenchmarkSingleAccess_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 只访问一个字段
		node, err := sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
			ConcurrentRead: true,
		}, "statuses", "0", "text")
		if err == nil && node.Exists() {
			_, _ = node.String()
		}
	}
}

// ============= 多次重复访问测试 =============

// 测试多次重复访问相同路径 - LoadAll 方式
func BenchmarkRepeatedAccess_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(bigJSONData)
		if err != nil {
			b.Fatal(err)
		}

		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 多次访问相同路径
		for j := 0; j < 50; j++ {
			node := root.GetByPath("statuses", "0", "text")
			if node.Exists() {
				_, _ = node.String()
			}
		}
	}
}

// 测试多次重复访问相同路径 - ConcurrentRead 方式
func BenchmarkRepeatedAccess_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 多次访问相同路径
		for j := 0; j < 50; j++ {
			node, err := sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "text")
			if err == nil && node.Exists() {
				_, _ = node.String()
			}
		}
	}
}

// ============= 重复访问所有字段测试 =============

// 测试重复访问所有字段 - LoadAll 方式 (Big JSON)
func BenchmarkRepeatedAccessAllFields_LoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := sonic.Get(bigJSONData)
		if err != nil {
			b.Fatal(err)
		}

		err = root.LoadAll()
		if err != nil {
			b.Fatal(err)
		}

		// 重复访问多个关键字段 20 次
		for j := 0; j < 20; j++ {
			// 访问第一个状态的所有主要字段
			root.GetByPath("statuses", "0", "id")
			root.GetByPath("statuses", "0", "text")
			root.GetByPath("statuses", "0", "created_at")
			root.GetByPath("statuses", "0", "source")
			root.GetByPath("statuses", "0", "truncated")
			root.GetByPath("statuses", "0", "in_reply_to_screen_name")

			// 访问用户信息字段
			root.GetByPath("statuses", "0", "user", "id")
			root.GetByPath("statuses", "0", "user", "name")
			root.GetByPath("statuses", "0", "user", "screen_name")
			root.GetByPath("statuses", "0", "user", "location")
			root.GetByPath("statuses", "0", "user", "description")
			root.GetByPath("statuses", "0", "user", "followers_count")
			root.GetByPath("statuses", "0", "user", "friends_count")
			root.GetByPath("statuses", "0", "user", "statuses_count")
			root.GetByPath("statuses", "0", "user", "verified")

			// 访问元数据字段
			root.GetByPath("statuses", "0", "metadata", "result_type")
			root.GetByPath("statuses", "0", "metadata", "iso_language_code")

			// 访问搜索元数据字段
			root.GetByPath("search_metadata", "completed_in")
			root.GetByPath("search_metadata", "max_id")
			root.GetByPath("search_metadata", "query")
			root.GetByPath("search_metadata", "count")
		}
	}
}

// 测试重复访问所有字段 - ConcurrentRead 方式 (Big JSON)
func BenchmarkRepeatedAccessAllFields_ConcurrentRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 重复访问多个关键字段 20 次
		for j := 0; j < 20; j++ {
			// 访问第一个状态的所有主要字段
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "id")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "text")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "created_at")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "source")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "truncated")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "in_reply_to_screen_name")

			// 访问用户信息字段
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "id")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "name")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "screen_name")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "location")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "description")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "followers_count")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "friends_count")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "statuses_count")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "user", "verified")

			// 访问元数据字段
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "metadata", "result_type")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "statuses", "0", "metadata", "iso_language_code")

			// 访问搜索元数据字段
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "search_metadata", "completed_in")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "search_metadata", "max_id")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "search_metadata", "query")
			sonic.GetWithOptions(bigJSONData, ast.SearchOptions{
				ConcurrentRead: true,
			}, "search_metadata", "count")
		}
	}
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
