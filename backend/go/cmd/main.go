package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/PusztaiMate/clipper-go-backend/clips"
	"github.com/PusztaiMate/clipper-go-backend/pkg/clippersrvc"
	"github.com/PusztaiMate/clipper-go-backend/pkg/server"
	"github.com/PusztaiMate/clipper-go-backend/pkg/utils"
	"google.golang.org/grpc"
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatalf("can not listen on address %s (maybe port is taken?)", lis.Addr())
	}
	s := grpc.NewServer()
	cs := clippersrvc.NewClipperService(logger, srcDir, clipDir)
	clips.RegisterClipsServer(s, server.NewClipServer(logger, cs))
	logger.Printf("starting server...")
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("server stopped: %v", err)
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
