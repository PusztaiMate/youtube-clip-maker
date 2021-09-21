package clips

func NewClipResponse(message string) *ClipsResponse {
	return &ClipsResponse{Message: message}
}
