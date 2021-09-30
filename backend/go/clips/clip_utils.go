package clips

import (
	"fmt"

	"github.com/PusztaiMate/clipper-go-backend/pkg/clippersrvc"
)

func NewClipResponse(message string) *ClipsResponse {
	return &ClipsResponse{Message: message}
}

func ToClip(c *ClipsRequest) *clippersrvc.ClipRequest {
	cr := &clippersrvc.ClipRequest{}

	cr.Url = c.Url
	for _, c := range c.GetClips() {
		cr.Clips = append(cr.Clips, clippersrvc.Clip{Player: c.Player, StartTime: c.StartTime, EndTime: c.EndTime})
	}

	return cr
}

func (c *Clip) ToString() string {
	return fmt.Sprintf("Clip(player: %s, start_time: %s, end_time: %s", c.GetPlayer(), c.GetStartTime(), c.GetEndTime())
}
