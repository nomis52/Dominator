package imaginator

import "time"

type BuildImageRequest struct {
	StreamName   string
	ExpiresIn    time.Duration
	GitBranch    string
	MaxSourceAge time.Duration
}

type BuildImageResponse struct {
	ImageName   string
	BuildLog    []byte
	ErrorString string
}
