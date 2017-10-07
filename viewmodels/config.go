package viewmodels

type Config struct {
	Title   string
	Heading string
}

func GetConfig() Config {
	return Config{
		Title:   "Config",
		Heading: "Application Config",
	}
}
