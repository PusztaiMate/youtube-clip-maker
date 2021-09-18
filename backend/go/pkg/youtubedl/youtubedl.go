package youtubedl

import (
	"log"
	"os/exec"
	"strings"
)

const (
	ARCHIVE = "archive.txt"
)

type YoutubeDlCommand struct {
	logger        *log.Logger
	url, outTempl string
}

func NewYoutubeDlCommand(logger *log.Logger, url, outTempl string) *YoutubeDlCommand {
	return &YoutubeDlCommand{logger, url, outTempl}
}

func (ydl YoutubeDlCommand) GetID() (string, error) {
	cmd := exec.Command("youtube-dl", "--get-id", ydl.url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	var videoId []byte
	_, err = stdout.Read(videoId)
	if err != nil {
		return "", err
	}

	return string(videoId), nil
}

func (ydl YoutubeDlCommand) Execute() error {
	cmd := exec.Command("youtube-dl", "-o", ydl.outTempl, "--download-archive", ARCHIVE, "-f", "best", ydl.url)
	ydl.logger.Printf("executing %s", strings.Join(cmd.Args, " "))
	return cmd.Run()
}
