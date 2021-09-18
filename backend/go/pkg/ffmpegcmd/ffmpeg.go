package ffmpegcmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type ffmpegCut struct {
	startTime, endTime string
}
type FfmpegCommand struct {
	logger                     *log.Logger
	cuts                       []ffmpegCut
	sourceVideoPath, outPrefix string
}

func NewFfmpegCommand(l *log.Logger, sourceVideo, outputNamePrefix string) *FfmpegCommand {
	cuts := make([]ffmpegCut, 0)
	return &FfmpegCommand{l, cuts, sourceVideo, outputNamePrefix}
}

func (f FfmpegCommand) Do() {
	var wg sync.WaitGroup

	for _, cut := range f.cuts {
		outName := f.createOutputName(cut.startTime, cut.endTime)
		cmd := createSingleCmd(f.sourceVideoPath, outName, cut.startTime, cut.endTime)

		wg.Add(1)
		go func(c *exec.Cmd) {
			defer wg.Done()
			f.logger.Printf("executing '%s'", strings.Join(c.Args, " "))
			err := c.Run()
			if err != nil {
				f.logger.Printf("error while issues command '%s': '%s'", strings.Join(c.Args, " "), err)
			}
		}(cmd)
	}
	f.logger.Println("Waiting for clips to be cut out...")
	wg.Wait()
	f.logger.Println("Cutting done!")
}

func (f *FfmpegCommand) AddCut(startTime, endTime string) {
	f.cuts = append(f.cuts, ffmpegCut{startTime, endTime})
}

func (f *FfmpegCommand) createOutputName(first, second string) string {
	return fmt.Sprintf("%s_%s_%s.mp4", f.outPrefix, first, second)
}

func createSingleCmd(input, output, startTime, endTime string) *exec.Cmd {
	return exec.Command("ffmpeg", "-ss", startTime, "-i", input, "-to", endTime, output)
}

func (f *FfmpegCommand) SetOutPrefix(prefix string) {
	f.outPrefix = prefix
}

func (f *FfmpegCommand) SetSourcePath(path string) {
	f.sourceVideoPath = path
}

func (f *FfmpegCommand) Reset() {
	f.sourceVideoPath = ""
	f.outPrefix = ""
	f.cuts = make([]ffmpegCut, 0)
}
