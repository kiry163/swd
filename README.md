# swd

一个 Go 敏感词检查库。分类由使用方自行定义，库只保留 `Word.Type` 字段，不内置固定分类。

## 安装

```bash
go get github.com/kiry163/swd
```

## 基本使用

```go
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
```

`words.txt` 支持一行一个词，也支持 `词,类型`：

```text
诈骗,risk
赌博,risk
自定义词,custom
```

空行和 `#` 开头的注释行会被跳过。单行最大长度为 1MiB；如果词为空，加载错误会包含行号。

## 多来源加载

`Load` 接收多个 loader，并把加载结果追加到当前词表后自动重建匹配器：

```go
err := eng.Load(
	context.Background(),
	swd.NewFileLoader("words.txt"),
	swd.NewReaderLoader(strings.NewReader("违规,policy\n")),
	swd.NewMemoryLoader([]swd.Word{{Text: "敏感词", Type: "custom"}}),
)
```

## 运行期更新

`AddWord`、`AddWords`、`RemoveWord`、`RemoveWords`、`Clear` 都会自动重建匹配器，适合运行过程中小范围增删词。
这些操作可以和 `FindAll`、`Find`、`Contains`、`Replace` 并发调用；更新成功后新匹配器会一次性发布。

```go
_ = eng.AddWord(swd.Word{Text: "新增词", Type: "custom"})
_ = eng.RemoveWord("旧词")
_ = eng.Clear()
```

如果需要整表替换，可以显式组合：

```go
_ = eng.Clear()
_ = eng.Load(context.Background(), swd.NewMemoryLoader(newWords))
```

## 匹配策略

`FindAll` 返回稳定的非重叠结果：

- 按原文位置从左到右返回。
- 同起点命中多个词时保留最长词。
- 被已选中长词覆盖的短词不会返回。
- 相同 `Text` 的词条后添加会覆盖先添加的 `Type` 和 `Meta`。

## 性能基准

项目包含构建、检测、查找、替换和运行期小范围更新的 benchmark：

```bash
go test -run '^$' -bench . -benchmem ./...
```
