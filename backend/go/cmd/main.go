package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PusztaiMate/clipper-go-backend/pkg/clippersrvc"
	"github.com/PusztaiMate/clipper-go-backend/pkg/restserver"
	"github.com/PusztaiMate/clipper-go-backend/pkg/utils"
)

const (
	ENV_CLIPS_DIR = "CLIPS_DIR"
	ENV_SRC_DIR   = "SRC_DIR"
	ENV_PORT      = "PORT"
)

func main() {
	logger := log.New(os.Stdout, "[CLIPPER SERVICE] ", log.LstdFlags)
	err := utils.CheckEnvVars(ENV_CLIPS_DIR, ENV_SRC_DIR, ENV_PORT)
	if err != nil {
		logger.Fatal(err)
	}

	srcDir, clipDir, err := createDirs()
	if err != nil {
		logger.Fatal(err)
	}
	port := os.Getenv(ENV_PORT)

	cs := clippersrvc.NewClipperService(logger, srcDir, clipDir)
	restServer := restserver.NewClipperRESTServer(logger, cs, "0.0.0.0", port)
	cancel := restServer.Run()
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals
	restServer.Shutdown()
}

func createDirs() (string, string, error) {
	srcDir, clipDir := os.Getenv(ENV_SRC_DIR), os.Getenv(ENV_CLIPS_DIR)

	err := utils.MakeDirs(srcDir, clipDir)
	if err != nil {
		return "", "", err
	}

	return srcDir, clipDir, nil
}
