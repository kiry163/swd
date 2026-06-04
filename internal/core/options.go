package core

type Config struct {
	IgnoreSymbol bool
	IgnoreWidth  bool
	IgnoreCase   bool
	Traditional  bool
	SimilarChar  bool
}

type Options = Config

type Option func(*Config)

func defaultConfig() Config {
	return Config{}
}

func WithIgnoreSymbol(v bool) Option { return func(c *Config) { c.IgnoreSymbol = v } }
func WithIgnoreWidth(v bool) Option  { return func(c *Config) { c.IgnoreWidth = v } }
func WithIgnoreCase(v bool) Option   { return func(c *Config) { c.IgnoreCase = v } }
func WithTraditional(v bool) Option  { return func(c *Config) { c.Traditional = v } }
func WithSimilarChar(v bool) Option  { return func(c *Config) { c.SimilarChar = v } }
