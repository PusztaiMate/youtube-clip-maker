package clippersrvc

import (
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PusztaiMate/clipper-go-backend/pkg/ffmpegcmd"
	"github.com/PusztaiMate/clipper-go-backend/pkg/utils"
	"github.com/PusztaiMate/clipper-go-backend/pkg/youtubedl"
)

type Clip struct {
	Player, StartTime, EndTime string
}
type ClipRequest struct {
	Url   string
	Clips []Clip
}

type ClipperService struct {
	logger               *log.Logger
	ytd                  *youtubedl.YoutubeDlCommand
	url, srcDir, clipDir string
}

func NewClipperService(l *log.Logger, srcDir, clipDir string) *ClipperService {
	ytd := youtubedl.NewYoutubeDlCommand(l, "", "")
	return &ClipperService{l, ytd, "", srcDir, clipDir}
}

func (cs *ClipperService) CreateClips(cr *ClipRequest) error {
	if cr == nil {
		return fmt.Errorf("the clip request is empty nil")
	}

	cs.setURL(cr.Url)
	videoId, err := cs.ytd.GetID()
	if err != nil {
		return err
	}

	ydlOut := path.Join(cs.srcDir, videoId)
	cs.ytd.SetOutTmpl(ydlOut)

	err = cs.ytd.Execute()
	if err != nil {
		cs.logger.Printf("could not download video (url: %s)", cr.Url)
		return err
	}

	//mostly because of unpredictable extension?
	sourceVideo, err := findSourceVideoInDir(cs.srcDir, videoId)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, c := range cr.Clips {
		wg.Add(1)
		go func(cc Clip) {
			defer wg.Done()
			playersDir := path.Join(cs.clipDir, cc.Player)
			err := utils.MakeDirs(playersDir)
			if err != nil {
				cs.logger.Print(err)
			}

			ffmpegOutTmpl := path.Join(playersDir, videoId)
			ffmpegCmd := ffmpegcmd.NewFfmpegCommand(cs.logger, sourceVideo, ffmpegOutTmpl)
			ffmpegCmd.AddCut(cc.StartTime, cc.EndTime)
			cs.logger.Printf("start cutting for %s", ffmpegCmd)
			ffmpegCmd.Do()
		}(c)
	}
	wg.Wait()

	return nil
}

func findSourceVideoInDir(dir string, videoId string) (string, error) {
	foundTwice := false
	var pathToSource string

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if !strings.Contains(path, videoId) {
			return nil
		}

		if pathToSource != "" {
			foundTwice = true
			return fmt.Errorf("%s found at least twice, aborting", videoId)
		}

		pathToSource = path
		return nil
	})

	if foundTwice {
		return "", fmt.Errorf("%s found at least twice, aborting", videoId)
	}
	return pathToSource, nil
}

func (cs *ClipperService) setURL(url string) error {
	cs.ytd.SetUrl(url)

	cs.url = url
	return nil
}
