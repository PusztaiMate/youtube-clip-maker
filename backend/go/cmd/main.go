package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/PusztaiMate/clipper-go-backend/clips"
	"google.golang.org/grpc"
)

const (
	ENV_CLIPS_DIR = "CLIPS_DIR"
	ENV_SRC_DIR   = "SRC_DIR"
	ENV_PORT      = "PORT"
)

type server struct {
	clips.UnimplementedClipServerServer

	logger *log.Logger
}

func (s *server) NewClip(ctx context.Context, ci *clips.ClipInfo) (*clips.Response, error) {
	s.logger.Printf("received %#v", ci)
	return &clips.Response{Content: "Hello"}, nil
}

func main() {
	logger := log.New(os.Stdout, "[CLIPPER SERVICE] ", log.LstdFlags)
	err := checkEnvVars()
	if err != nil {
		logger.Fatal(err)
	}

	port := os.Getenv(ENV_PORT)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatalf("can not listen on address %s (maybe port is taken?)", lis.Addr())
	}
	s := grpc.NewServer()
	clips.RegisterClipServerServer(s, &server{logger: logger})
	logger.Printf("starting server...")
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("server stopped: %v", err)
	}
}

func checkEnvVars() error {
	var missing []string
	envVars := []string{ENV_CLIPS_DIR, ENV_PORT, ENV_SRC_DIR}

	for _, v := range envVars {
		if _, ok := os.LookupEnv(v); !ok {
			missing = append(missing, v)
		}
	}

	if len(missing) != 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}
