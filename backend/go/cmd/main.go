package main

import (
	"log"
	"os"

	"github.com/PusztaiMate/clipper-go-backend/pkg/clippersrvc"
	"github.com/PusztaiMate/clipper-go-backend/pkg/restserver"
	"github.com/PusztaiMate/clipper-go-backend/pkg/utils"
)

const (
	ENV_CLIPS_DIR = "CLIPS_DIR"
	ENV_SRC_DIR   = "SRC_DIR"
	ENV_REST_PORT = "REST_PORT"
	ENV_GRPC_PORT = "GRPC_PORT"
)

func main() {
	logger := log.New(os.Stdout, "[CLIPPER SERVICE] ", log.LstdFlags)
	err := utils.CheckEnvVars(ENV_CLIPS_DIR, ENV_SRC_DIR, ENV_REST_PORT, ENV_GRPC_PORT)
	if err != nil {
		logger.Fatal(err)
	}

	srcDir, clipDir, err := createDirs()
	if err != nil {
		logger.Fatal(err)
	}
	restPort := os.Getenv(ENV_REST_PORT)

	cs := clippersrvc.NewClipperService(logger, srcDir, clipDir)

	logger.Printf("starting rest server at port %s", restPort)
	restServer := restserver.NewClipperRESTServer(logger, cs, "0.0.0.0", restPort)
	err = restServer.Run()
	if err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}

func createDirs() (string, string, error) {
	srcDir, clipDir := os.Getenv(ENV_SRC_DIR), os.Getenv(ENV_CLIPS_DIR)

	err := utils.MakeDirs(srcDir, clipDir)
	if err != nil {
		return "", "", err
	}

	return srcDir, clipDir, nil
}
