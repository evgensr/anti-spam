package main

import (
	"fmt"

	"github.com/evgensr/anti-spam/internal/app"
	"github.com/evgensr/anti-spam/internal/config"
	"github.com/evgensr/anti-spam/internal/logging"
	"github.com/evgensr/anti-spam/internal/telegram"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	l := logging.NewLogger(cfg)
	defer func() {
		_ = l.Sync()
	}()

	t, err := telegram.NewTelegram(cfg)
	if err != nil {
		fmt.Println("pls see log file")
		l.Fatal("error init telegram", zap.Error(err))
	}

	a := app.NewApp(cfg, l, t)
	a.Run()

}
