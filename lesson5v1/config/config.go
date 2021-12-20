package config

//Config - структура для конфигурации
type Config struct {
	MaxDepth   int
	MaxResults int
	MaxErrors  int
	Url        string
	Timeout    int //in seconds
}

func NewConfig(maxDepth int, maxResults int, maxErrors int, url string, timeout int) *Config {
	return &Config{
		MaxDepth:   maxDepth,
		MaxResults: maxResults,
		MaxErrors:  maxErrors,
		Url:        url,
		Timeout:    timeout,
	}
}
