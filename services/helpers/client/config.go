package client

type Config struct {
	Insecure     bool
	Proxy        string
	SmartlingURL string
}

func NewConfig(insecure bool, proxy string, smartlingURL string) Config {
	return Config{
		Insecure:     insecure,
		Proxy:        proxy,
		SmartlingURL: smartlingURL,
	}
}
