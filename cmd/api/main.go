package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/maybemaby/sk-connect/api"
)

type Args struct {
	Port string
}

func argParse() Args {
	var args Args
	flag.StringVar(&args.Port, "port", "8000", "port to listen on")
	flag.Parse()
	return args
}

func main() {

	// OS signal handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	args := argParse()

	appEnv := os.Getenv("APP_ENV")
	allowedHosts := os.Getenv("ALLOWED_HOSTS")

	hosts := strings.Split(allowedHosts, ",")
	isDebug := appEnv == "development"

	var logLevel slog.Level

	if isDebug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	cfg := api.ServerConfig{
		Port:         args.Port,
		LogLevel:     logLevel,
		AllowedHosts: hosts,
	}

	server := api.NewServer(cfg)

	go func() {
		err := server.Start()

		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-osSignals
}
