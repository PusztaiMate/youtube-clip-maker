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
	var wgExec sync.WaitGroup
	errChan := make(chan error)

	for _, cut := range f.cuts {
		outName := f.createOutputName(cut.startTime, cut.endTime)
		cmd := createSingleCmd(f.sourceVideoPath, outName, cut.startTime, cut.endTime)

		wgExec.Add(1)
		go func(c *exec.Cmd) {
			defer wgExec.Done()
			f.logger.Printf("executing '%s'", strings.Join(c.Args, " "))
			out, err := c.CombinedOutput()
			if err != nil {
				errChan <- fmt.Errorf("error while issuing command '%s': '%s' (%s)", c, err, out)
			}
		}(cmd)
	}
	go func() {
		wgExec.Wait()
		close(errChan)
	}()

	f.logger.Println("Waiting for clips to be cut out...")
	for err := range errChan {
		f.logger.Printf("%s", err)
	}
	f.logger.Println("Cutting done!")
}

func (f *FfmpegCommand) AddCut(startTime, endTime string) {
	f.cuts = append(f.cuts, ffmpegCut{startTime, endTime})
}

func (f *FfmpegCommand) createOutputName(first, second string) string {
	return fmt.Sprintf("%s_%s_%s.mp4", f.outPrefix, first, second)
}

func createSingleCmd(input, output, startTime, endTime string) *exec.Cmd {
	return exec.Command("ffmpeg", "-ss", startTime, "-i", input, "-to", endTime, "-y", replaceSpaceWithUnderscore((output)))
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

func replaceSpaceWithUnderscore(s string) string {
	return strings.Replace(s, " ", "_", 0)
}
