package swd

import (
	"io"

	internalengine "github.com/kiry163/swd/internal/engine"
)

type (
	Word   = internalengine.Word
	Match  = internalengine.Match
	Config = internalengine.Config
	Option = internalengine.Option
	Loader = internalengine.Loader
	Engine = internalengine.Engine
)

func New(opts ...Option) *Engine { return internalengine.New(opts...) }

func NewMemoryLoader(words []Word) Loader { return internalengine.NewMemoryLoader(words) }
func NewFileLoader(path string) Loader    { return internalengine.NewFileLoader(path) }
func NewReaderLoader(r io.Reader) Loader  { return internalengine.NewReaderLoader(r) }

func WithIgnoreSymbol(v bool) Option { return internalengine.WithIgnoreSymbol(v) }
func WithIgnoreWidth(v bool) Option  { return internalengine.WithIgnoreWidth(v) }
func WithIgnoreCase(v bool) Option   { return internalengine.WithIgnoreCase(v) }
func WithTraditional(v bool) Option  { return internalengine.WithTraditional(v) }
func WithPinyin(v bool) Option       { return internalengine.WithPinyin(v) }
func WithSimilarChar(v bool) Option  { return internalengine.WithSimilarChar(v) }
func WithHomophone(v bool) Option    { return internalengine.WithHomophone(v) }
