package swd

import internalcore "github.com/kiry163/swd/internal/core"

func New(opts ...Option) *Engine { return internalcore.New(opts...) }

func NewWithOptions(options Options) *Engine { return internalcore.NewWithOptions(options) }
