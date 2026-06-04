package main

import (
	"context"
	"fmt"
	"github.com/kiry163/swd"
)

func main() {
	eng := swd.New(
		swd.WithIgnoreSymbol(true),
		swd.WithIgnoreCase(true),
	)

	err := eng.Load(
		context.Background(),
		swd.NewFileLoader("words.txt"),
		swd.NewMemoryLoader([]swd.Word{
			{Text: "测试", Type: "custom"},
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(eng.Contains("这是一段测@试文本"))
	fmt.Println(eng.Replace("这是一段测试文本", "*"))
}
