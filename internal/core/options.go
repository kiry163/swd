package core

type Config struct {
	IgnoreSymbol bool
	IgnoreWidth  bool
	IgnoreCase   bool
	Traditional  bool
	SimilarChar  bool
}

type Options = Config

func defaultConfig() Config {
	return Config{}
}
