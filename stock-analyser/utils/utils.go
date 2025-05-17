package utils

type AppConfig struct {
	ServerHost string
	ServerPort string
}

var Config *AppConfig

func UpdateVariables() {
	Config = &AppConfig{
		ServerHost: "localhost",
		ServerPort: "8080",
	}
}
