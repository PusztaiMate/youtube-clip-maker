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
	clips.UnimplementedClipsServer

	logger *log.Logger
}

func (s *server) NewClip(c context.Context, cr *clips.ClipsRequest) (*clips.ClipsResponse, error) {
	s.logger.Printf("url in the incoming request: %s", cr.Url)
	s.logger.Printf("received %d clips", len(cr.Clips))
	for _, clip := range cr.Clips {
		s.logger.Printf("received clip: %v", clip)
	}
	return &clips.ClipsResponse{Message: "success"}, nil
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
	clips.RegisterClipsServer(s, &server{logger: logger})
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
