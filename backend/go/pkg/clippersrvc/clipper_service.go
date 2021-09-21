package clippersrvc

import (
	"fmt"
	"log"
	"path"
	"sync"

	"github.com/PusztaiMate/clipper-go-backend/clips"
	"github.com/PusztaiMate/clipper-go-backend/pkg/ffmpegcmd"
	"github.com/PusztaiMate/clipper-go-backend/pkg/utils"
	"github.com/PusztaiMate/clipper-go-backend/pkg/youtubedl"
)

type ClipperService struct {
	logger               *log.Logger
	ytd                  *youtubedl.YoutubeDlCommand
	url, srcDir, clipDir string
}

func NewClipperService(l *log.Logger, srcDir, clipDir string) *ClipperService {
	ytd := youtubedl.NewYoutubeDlCommand(l, "", "")
	return &ClipperService{l, ytd, "", srcDir, clipDir}
}

func (cs *ClipperService) CreateClips(cr *clips.ClipsRequest) error {
	if cr == nil {
		return fmt.Errorf("the clip request is empty nil")
	}

	cs.setURL(cr.Url)
	videoId, err := cs.ytd.GetID()
	if err != nil {
		return err
	}

	ydlOut := path.Join(cs.srcDir, fmt.Sprintf("%s.mp4", videoId))
	cs.ytd.SetOutTmpl(ydlOut)

	err = cs.ytd.Execute()
	if err != nil {
		cs.logger.Printf("could not download video (url: %s)", cr.Url)
		return err
	}

	var wg sync.WaitGroup
	for _, c := range cr.GetClips() {
		wg.Add(1)
		go func(cc *clips.Clip) {
			defer wg.Done()
			playersDir := path.Join(cs.clipDir, cc.Player)
			err := utils.MakeDirs(playersDir)
			if err != nil {
				cs.logger.Print(err)
			}

			ffmpegOutTmpl := path.Join(playersDir, videoId)
			ffmpegCmd := ffmpegcmd.NewFfmpegCommand(cs.logger, ydlOut, ffmpegOutTmpl)
			ffmpegCmd.AddCut(cc.StartTime, cc.EndTime)
			ffmpegCmd.Do()
		}(c)
	}
	wg.Wait()

	return nil
}

func (cs *ClipperService) setURL(url string) error {
	cs.ytd.SetUrl(url)

	cs.url = url
	return nil
}
