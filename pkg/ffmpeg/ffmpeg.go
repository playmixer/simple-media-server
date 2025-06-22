package ffmpeg

import "os/exec"

type FFMpeg struct {
}

func New() (*FFMpeg, error) {
	return &FFMpeg{}, nil
}

func (f *FFMpeg) AviToMP4(from, to string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", from, to)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}
