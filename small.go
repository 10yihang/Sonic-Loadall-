package test

// ffjson: skip
// easyjson:skip
type Book struct {
	BookId  int       `json:"id"`
	BookIds []int     `json:"ids"`
	Title   string    `json:"title"`
	Titles  []string  `json:"titles"`
	Price   float64   `json:"price"`
	Prices  []float64 `json:"prices"`
	Hot     bool      `json:"hot"`
	Hots    []bool    `json:"hots"`
	Author  Author    `json:"author"`
	Authors []Author  `json:"authors"`
	Weights []int     `json:"weights"`
}

// ffjson: skip
// easyjson:skip
type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Male bool   `json:"male"`
}

var book = Book{
	BookId:  12125925,
	BookIds: []int{-2147483648, 2147483647},
	Title:   "未来简史-从智人到智神",
	Titles:  []string{"hello", "world"},
	Price:   40.8,
	Prices:  []float64{-0.1, 0.1},
	Hot:     true,
	Hots:    []bool{true, true, true},
	Author:  author,
	Authors: []Author{author, author, author},
	Weights: nil,
}

var author = Author{
	Name: "json",
	Age:  99,
	Male: true,
}