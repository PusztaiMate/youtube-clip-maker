package main

import (
	"log"
	"net/http"
	"os/exec"
)

type FfmpegCommand struct {
	logger      *log.Logger
	commandArgs []string
}

func NewFfmpegCommand(l *log.Logger) *FfmpegCommand {
	commandArgs := make([]string, 0)
	return &FfmpegCommand{l, commandArgs}
}

func (f FfmpegCommand) Do() {
	cmd := exec.Command("ffmpeg", f.commandArgs...)
	err := cmd.Run()
	if err != nil {
		log.Printf("")
	}
}

func main() {
	server := http.Server{}

	server.ListenAndServe()
}
