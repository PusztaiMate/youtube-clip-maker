package clippersrvc

import (
	"log"

	"github.com/PusztaiMate/clipper-go-backend/pkg/ffmpegcmd"
	"github.com/PusztaiMate/clipper-go-backend/pkg/youtubedl"
)

type ClipperService struct {
	logger *log.Logger
	ffmpeg *ffmpegcmd.FfmpegCommand
	ytd    *youtubedl.YoutubeDlCommand
	url    string
}

func NewClipperService(l *log.Logger) *ClipperService {
	ffmpeg := ffmpegcmd.NewFfmpegCommand(l, "", "")
	ytd := youtubedl.NewYoutubeDlCommand(l, "", "")
	return &ClipperService{l, ffmpeg, ytd, ""}
}

func (cs *ClipperService) AddClip(startTime, endTime string) {
	cs.ffmpeg.AddCut(startTime, endTime)
}

func (cs *ClipperService) SetURL(url string) error {
	videoId, err := cs.ytd.GetID()
	if err != nil {
		return err
	}

	cs.ffmpeg.SetOutPrefix(videoId)
	cs.url = url
	return nil
}

func (cs *ClipperService) CreateClips() {
	cs.ffmpeg.Do()
	cs.ffmpeg.Reset()
}
