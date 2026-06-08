package swd

import (
	"io"

	internalcore "github.com/kiry163/swd/internal/core"
)

func NewMemoryLoader(words []Word) Loader { return internalcore.NewMemoryLoader(words) }
func NewFileLoader(path string) Loader    { return internalcore.NewFileLoader(path) }
func NewReaderLoader(r io.Reader) Loader  { return internalcore.NewReaderLoader(r) }
func NewStringLoader(text string) Loader  { return internalcore.NewStringLoader(text) }
