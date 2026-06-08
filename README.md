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
	eng := swd.New(swd.Options{
		IgnoreSymbol: true,
		IgnoreCase:   true,
	})

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

`Replace` 会把命中范围内的每个原文字符替换为一次完整 `mask`。例如敏感词长度为 2 且 `mask` 为 `[x]` 时，替换结果会使用 `[x][x]`。

## 归一化选项

`Options.Traditional` 会使用 OpenCC 的 `t2s` 规则把词库和输入文本统一归一到简体后再匹配。这个选项只处理繁体到简体转换，不包含台湾、香港等地区词汇转换；例如 `軟體` 会归一为 `软体`，不会归一为 `软件`。

`Options.IgnoreWidth` 会使用 Unicode width folding 归一化宽窄字符，例如全角 ASCII、全角空格和半角片假名。词库里的 `ガス` 可以匹配输入文本中的 `ｶﾞｽ`，命中结果仍返回原文位置。

`Options.IgnoreSymbol` 会忽略空白、标点、符号、控制字符和零宽格式字符，适合处理 `敏​感‍词` 这类插入干扰字符的文本。

`Options.SimilarChar` 只做低误报的常见字符混淆归一化，例如 `b@d` 匹配 `bad`、`t3st` 匹配 `test`。库不提供拼音和同音字匹配，以避免不可控的误报。

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

如果需要整表替换，使用 `ReplaceWords`。它会先完整校验并构建新匹配器，成功后再一次性发布；如果失败，旧词表会继续生效。

```go
_ = eng.ReplaceWords(newWords)
```

启用归一化选项后，如果某个词归一化后变为空字符串，例如开启 `IgnoreSymbol` 后加载纯符号词，加载和更新会返回错误，避免词条静默失效。

## 导出词表

`Words` 会返回当前词表快照；`Export` 和 `ExportFile` 可以把当前词表导出成 loader 可再次读取的简易格式：

```go
words := eng.Words()
_ = words

_ = eng.ExportFile(context.Background(), "words.txt")
```

导出格式和加载格式一致：无类型写成 `词`，有类型写成 `词,类型`。`ExportFile` 会覆盖目标文件。当前简易格式不支持词或类型中包含逗号、回车或换行；遇到这类词条时导出会返回错误。

## 匹配策略

`FindAll` 返回稳定的非重叠结果：

- 按原文位置从左到右返回。
- 同起点命中多个词时保留最长词。
- 被已选中长词覆盖的短词不会返回。
- `Load`、`AddWords`、`ReplaceWords` 遇到相同 `Text` 的词条时，后出现的词条会覆盖先出现的 `Type`，并保留第一次出现的位置。

## 性能基准

项目包含构建、检测、查找、替换和运行期小范围更新的 benchmark：

```bash
go test -run '^$' -bench . -benchmem ./...
```
