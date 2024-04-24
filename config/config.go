package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port string `json:"server_port"` // Port on which the webserver listens.
	Path string `json:"path"`        // path to the files to load
}

const defaultSrvPort = "8080"
const defaultPath = "/usr/items"

var CONFIG Config

func InitConfig() {
	CONFIG = Config{}
	CONFIG.Port = defaultSrvPort
	CONFIG.Path = defaultPath

	if srvPort, ok := os.LookupEnv("SERVER_PORT"); ok {
		CONFIG.Port = srvPort
	}

	if path, ok := os.LookupEnv("FILE_PATH"); ok {
		CONFIG.Path = path
	}

	PrintSettings()
}

func PrintSettings() {
	slog.Info("Settings", "port", CONFIG.Port, "path", CONFIG.Path)
}

func InitLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(logger)
}
