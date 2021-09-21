package server

import (
	"context"
	"fmt"
	"log"

	"github.com/PusztaiMate/clipper-go-backend/clips"
	"github.com/PusztaiMate/clipper-go-backend/pkg/clippersrvc"
)

type ClipServer struct {
	clips.UnimplementedClipsServer

	logger  *log.Logger
	clipper *clippersrvc.ClipperService
}

func NewClipServer(logger *log.Logger, clipper *clippersrvc.ClipperService) *ClipServer {
	return &ClipServer{logger: logger, clipper: clipper}
}

func (s *ClipServer) NewClip(c context.Context, cr *clips.ClipsRequest) (*clips.ClipsResponse, error) {
	err := s.clipper.CreateClips(cr)
	if err != nil {
		return clips.NewClipResponse(fmt.Sprintf("error when processing clips: %s", err)), err
	}

	return clips.NewClipResponse("success, clips will be created"), nil
}
