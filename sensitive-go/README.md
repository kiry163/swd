# Sensitive-Go

A high-performance, feature-rich sensitive word detection library for Go.

[![Go Version](https://img.shields.io/github/go-mod/go-version/Karrecy/sensitive-go?style=flat-square)](https://golang.org)
[![License](https://img.shields.io/github/license/Karrecy/sensitive-go?style=flat-square)](LICENSE)
[![Stars](https://img.shields.io/github/stars/Karrecy/sensitive-go?style=flat-square)](https://github.com/Karrecy/sensitive-go/stargazers)
[![Last Commit](https://img.shields.io/github/last-commit/Karrecy/sensitive-go?style=flat-square)](https://github.com/Karrecy/sensitive-go/commits/main)


[中文文档](README_zh.md)

## Features

- 🚀 **High Performance**: DFA and Aho-Corasick algorithms with auto-selection
- 🔧 **Variant Detection**: Pinyin, traditional Chinese, symbol filtering, and similar characters
- 🎯 **Flexible Matching**: Case-insensitive and whitelist support
- 📦 **Multiple Loaders**: File, HTTP, and memory sources for both blacklist and whitelist
- 🔄 **Auto Reload**: File monitoring with automatic dictionary updates
- 🔒 **Thread-Safe**: Safe for concurrent use
- 📦 **Zero Dependencies**: Core library with no external dependencies

## Installation

```bash
go get github.com/Karrecy/sensitive-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/Karrecy/sensitive-go"
)

func main() {
    // Use built-in dictionary (recommended for quick start)
    detector, _ := gosensitive.New().
        LoadBuiltin().  // Load built-in default dictionary
        Build()

    // Or load from custom sources
    detector, _ = gosensitive.New().
        LoadMemory([]string{"badword", "spam"}).
        Build()

    // Check if text contains sensitive words
    if detector.Contains("This is a badword") {
        fmt.Println("Sensitive word detected!")
    }

    // Find all matches
    matches := detector.Find("badword and spam")
    for _, match := range matches {
        fmt.Printf("Found: %s at [%d:%d]\n", match.Word, match.Start, match.End)
    }

    // Replace sensitive words
    filtered := detector.Filter("This badword is spam")
    fmt.Println(filtered) // Output: This ******* is ****
}
```

## Core Features

### 1. Algorithm Selection

```go
// Auto-select (DFA for <5000 words, AC for ≥5000)
detector := gosensitive.New().
    UseAlgorithm(gosensitive.AlgorithmAuto).
    LoadFile("words.txt").
    Build()

// Explicit selection
detector := gosensitive.New().
    UseAlgorithm(gosensitive.AlgorithmDFA).  // or AlgorithmAC
    LoadFile("words.txt").
    Build()
```

### 2. Case-Insensitive Matching

```go
detector := gosensitive.New().
    LoadMemory([]string{"Test", "Example"}).
    SetCaseSensitive(false).  // Case-insensitive
    Build()

// Matches: "test", "TEST", "Test", "tEsT"
fmt.Println(detector.Contains("this is a TEST"))  // true
```

### 3. Variant Detection

```go
detector := gosensitive.New().
    LoadMemory([]string{"测试"}).
    EnableSymbol().       // Remove symbols: "测@试" → "测试"
    EnableTraditional().  // Simplified/Traditional: "測試" → "测试"
    EnableSimilarChar().  // Similar chars: "测st" → "测试"
    EnablePinyin().       // Pinyin: "ceshi" → "测试"
    Build()

// Detects variants
detector.Contains("测@试")    // true (symbol removed)
detector.Contains("測試")     // true (traditional)
detector.Contains("ce shi")   // true (pinyin)
```

### 4. Whitelist Support

```go
// From memory
detector := gosensitive.New().
    LoadMemory([]string{"test", "example", "sensitive"}).
    AddWhitelist("test", "example").  // Exclude these
    Build()

// From file
detector := gosensitive.New().
    LoadFile("blacklist.txt").
    LoadWhitelistFile("whitelist.txt").  // Load from file
    Build()

// Multiple sources
detector := gosensitive.New().
    LoadFile("words.txt").
    LoadWhitelistFile("whitelist1.txt").
    LoadWhitelistHTTP("https://example.com/whitelist.txt").
    AddWhitelist("temporary").  // Add more
    Build()
```

### 5. Multiple Loading Sources

```go
// Built-in dictionary (embedded in binary)
detector := gosensitive.New().
    LoadBuiltin().  // Load built-in default dictionary
    Build()

// Multiple sources
detector := gosensitive.New().
    LoadBuiltin().                            // Built-in dictionary
    LoadFile("local_words.txt").              // Local file
    LoadHTTP("https://cdn.com/words.txt").    // Remote HTTP
    LoadMemory([]string{"extra1", "extra2"}). // Memory
    Build()
```

### 6. File Monitoring & Auto Reload

```go
opts := gosensitive.DefaultOptions()
opts.WatchFile = true
opts.WatchInterval = time.Second * 30  // Check every 30s

detector, _ := gosensitive.New().
    LoadFile("words.txt").
    SetOptions(opts).
    Build()

// File changes are automatically detected and reloaded
defer detector.Close()  // Stop watchers
```

### 7. Category & Level Filtering

```go
words := []dict.Word{
    {Text: "politics", Category: dict.CategoryPolitical, Level: dict.LevelHigh},
    {Text: "spam", Category: dict.CategoryAd, Level: dict.LevelLow},
}

opts := gosensitive.DefaultOptions()
opts.Categories = []Category{CategoryPolitical}  // Only political
opts.MinLevel = LevelHigh                        // Only high level

detector := gosensitive.New().
    LoadWords(words).
    SetOptions(opts).
    Build()
```

### 8. Custom Options

```go
opts := gosensitive.DefaultOptions()
opts.ReplaceChar = '█'
opts.MaxMatchCount = 10
opts.CaseSensitive = false

detector := gosensitive.New().
    LoadMemory([]string{"word1", "word2"}).
    SetOptions(opts).
    Build()
```

## Whitelist File Format

**Plain Text (whitelist.txt)**:
```text
test
example
# Comments are ignored
normal_word
```

**JSON (whitelist.json)**:
```json
[
  {"text": "test", "category": 0, "level": 0},
  {"text": "example", "category": 0, "level": 0}
]
```
