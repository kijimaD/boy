package main

import (
	"errors"
	"log"
	"os"

	"github.com/kijimaD/goboy/pkg/cartridge"
	"github.com/kijimaD/goboy/pkg/logger"
	"github.com/kijimaD/goboy/pkg/utils"
)

// go run main.go roms/helloworld/hello.gb

func main() {
	level := "Debug"
	if os.Getenv("LEVEL") != "" {
		level = os.Getenv("LEVEL")
	}
	l := logger.NewLogger(logger.LogLevel(level))
	if len(os.Args) != 2 {
		log.Fatalf("ERROR: %v", errors.New("Please specify the ROM"))
	}
	file := os.Args[1]
	log.Println(file)
	buf, err := utils.LoadROM(file)
	if err != nil {
		log.Fatalf("ERROR: %v", errors.New("Failed to load ROM"))
	}
	cart, err := cartridge.NewCartridge(buf)
	if err != nil {
		log.Fatalf("ERROR: %v", errors.New("Failed to create cartridge"))
	}

	l.Info(cart, "success load")
}
