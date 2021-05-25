package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/app/manager"
)

// LC - logging's category
const LC = "MAIN"

func main() {
	configInst := new(config.Config).Create()
	cfg := configInst.GetConfig()

	logCfg := logger.Cfg{
		AppName: cfg.Service.Name,
		LogPath: cfg.Logger.Path,
		Level:   cfg.Logger.Level,
	}

	logger.SetConfig(logCfg)

	logger.InfoJ(LC, "STARTING SERVICE")

	logger.DebugJ(LC, fmt.Sprint("Config: ", spew.Sdump(cfg)))

	managerInst := new(manager.Manager).Create(cfg)
	go managerInst.Run()

	<-managerInst.GetChCritError()

	logger.InfoJ(LC, "STOPING SERVICE")
}
