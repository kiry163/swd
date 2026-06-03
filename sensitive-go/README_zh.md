# Sensitive-Go

一个高性能、功能丰富的 Go 语言敏感词检测库。

[![Go Version](https://img.shields.io/github/go-mod/go-version/Karrecy/sensitive-go?style=flat-square)](https://golang.org)
[![License](https://img.shields.io/github/license/Karrecy/sensitive-go?style=flat-square)](LICENSE)
[![Stars](https://img.shields.io/github/stars/Karrecy/sensitive-go?style=flat-square)](https://github.com/Karrecy/sensitive-go/stargazers)
[![Last Commit](https://img.shields.io/github/last-commit/Karrecy/sensitive-go?style=flat-square)](https://github.com/Karrecy/sensitive-go/commits/main)


[English](README.md)

## 特性

- 🚀 **高性能**: DFA 和 Aho-Corasick 算法，自动选择最优方案
- 🔧 **变体检测**: 拼音、繁简体、符号干扰、形近字检测
- 🎯 **灵活匹配**: 大小写不敏感、白名单支持
- 📦 **多种加载方式**: 黑名单和白名单均支持文件、HTTP、内存加载
- 🔄 **自动重载**: 文件监控，自动更新词库
- 🔒 **线程安全**: 支持高并发使用
- 📦 **零依赖**: 核心库无外部依赖

## 安装

```bash
go get github.com/Karrecy/sensitive-go
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/Karrecy/sensitive-go"
)

func main() {
    // 使用内置词库（推荐快速开始）
    detector, _ := gosensitive.New().
        LoadBuiltin().  // 加载内置默认词库
        Build()

    // 或从自定义来源加载
    detector, _ = gosensitive.New().
        LoadMemory([]string{"敏感词", "测试"}).
        Build()

    // 检查是否包含敏感词
    if detector.Contains("这是一个敏感词") {
        fmt.Println("检测到敏感词！")
    }

    // 查找所有敏感词
    matches := detector.Find("敏感词和测试")
    for _, match := range matches {
        fmt.Printf("发现: %s 位置 [%d:%d]\n", match.Word, match.Start, match.End)
    }

    // 替换敏感词
    filtered := detector.Filter("这个敏感词需要过滤")
    fmt.Println(filtered) // 输出: 这个***需要过滤
}
```

## 核心功能

### 1. 算法选择

```go
// 自动选择（词库<5000用DFA，≥5000用AC）
detector := gosensitive.New().
    UseAlgorithm(gosensitive.AlgorithmAuto).
    LoadFile("words.txt").
    Build()

// 显式指定
detector := gosensitive.New().
    UseAlgorithm(gosensitive.AlgorithmDFA).  // 或 AlgorithmAC
    LoadFile("words.txt").
    Build()
```

### 2. 大小写不敏感匹配

```go
detector := gosensitive.New().
    LoadMemory([]string{"测试", "Test"}).
    SetCaseSensitive(false).  // 不区分大小写
    Build()

// 能匹配 "test", "TEST", "Test", "tEsT"
fmt.Println(detector.Contains("这是一个TEST"))  // true
```

### 3. 变体检测

```go
detector := gosensitive.New().
    LoadMemory([]string{"测试"}).
    EnableSymbol().       // 去除符号: "测@试" → "测试"
    EnableTraditional().  // 繁简转换: "測試" → "测试"
    EnableSimilarChar().  // 形近字: "测st" → "测试"
    EnablePinyin().       // 拼音: "ceshi" → "测试"
    Build()

// 检测变体
detector.Contains("测@试")    // true (去除符号)
detector.Contains("測試")     // true (繁体)
detector.Contains("ce shi")   // true (拼音)
```

### 4. 白名单支持

```go
// 从内存加载
detector := gosensitive.New().
    LoadMemory([]string{"测试", "示例", "敏感"}).
    AddWhitelist("测试", "示例").  // 排除这些词
    Build()

// 从文件加载
detector := gosensitive.New().
    LoadFile("blacklist.txt").
    LoadWhitelistFile("whitelist.txt").  // 从文件加载白名单
    Build()

// 多种来源
detector := gosensitive.New().
    LoadFile("words.txt").
    LoadWhitelistFile("whitelist1.txt").
    LoadWhitelistHTTP("https://example.com/whitelist.txt").
    AddWhitelist("临时豁免").  // 添加更多
    Build()
```

### 5. 多种加载方式

```go
// 内置词库（嵌入在二进制中）
detector := gosensitive.New().
    LoadBuiltin().  // 加载内置默认词库
    Build()

// 多种来源组合
detector := gosensitive.New().
    LoadBuiltin().                            // 内置词库
    LoadFile("local_words.txt").              // 本地文件
    LoadHTTP("https://cdn.com/words.txt").    // 远程HTTP
    LoadMemory([]string{"额外1", "额外2"}).   // 内存
    Build()
```

### 6. 文件监控和自动重载

```go
opts := gosensitive.DefaultOptions()
opts.WatchFile = true
opts.WatchInterval = time.Second * 30  // 每30秒检查一次

detector, _ := gosensitive.New().
    LoadFile("words.txt").
    SetOptions(opts).
    Build()

// 文件变化会自动检测并重载
defer detector.Close()  // 停止监控
```

### 7. 分类和等级过滤

```go
words := []dict.Word{
    {Text: "政治词", Category: dict.CategoryPolitical, Level: dict.LevelHigh},
    {Text: "广告词", Category: dict.CategoryAd, Level: dict.LevelLow},
}

opts := gosensitive.DefaultOptions()
opts.Categories = []Category{CategoryPolitical}  // 只检测政治类
opts.MinLevel = LevelHigh                        // 只检测高级别

detector := gosensitive.New().
    LoadWords(words).
    SetOptions(opts).
    Build()
```

### 8. 自定义选项

```go
opts := gosensitive.DefaultOptions()
opts.ReplaceChar = '█'
opts.MaxMatchCount = 10
opts.CaseSensitive = false

detector := gosensitive.New().
    LoadMemory([]string{"甲", "乙"}).
    SetOptions(opts).
    Build()
```

## 白名单文件格式

**纯文本格式(whitelist.txt)**:
```text
测试
示例
# 注释会被忽略
正常词汇
```

**JSON格式 (whitelist.json)**:
```json
[
  {"text": "测试", "category": 0, "level": 0},
  {"text": "示例", "category": 0, "level": 0}
]
```
