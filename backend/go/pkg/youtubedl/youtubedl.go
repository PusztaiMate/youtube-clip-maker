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
	out, err := exec.Command("youtube-dl", "--get-id", ydl.url).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (ydl YoutubeDlCommand) Execute() error {
	cmd := exec.Command("youtube-dl", "-o", ydl.outTempl, "--download-archive", ARCHIVE, ydl.url)
	ydl.logger.Printf("executing %s", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func (ydl *YoutubeDlCommand) SetOutTmpl(tmpl string) {
	ydl.outTempl = tmpl
}

func (ydl *YoutubeDlCommand) SetUrl(url string) {
	ydl.url = url
}

func (ydl *YoutubeDlCommand) Reset() {
	ydl.url = ""
	ydl.outTempl = ""
}
