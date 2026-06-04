package swd

import internalcore "github.com/kiry163/swd/internal/core"

func WithIgnoreSymbol(v bool) Option { return internalcore.WithIgnoreSymbol(v) }
func WithIgnoreWidth(v bool) Option  { return internalcore.WithIgnoreWidth(v) }
func WithIgnoreCase(v bool) Option   { return internalcore.WithIgnoreCase(v) }
func WithTraditional(v bool) Option  { return internalcore.WithTraditional(v) }
func WithSimilarChar(v bool) Option  { return internalcore.WithSimilarChar(v) }
