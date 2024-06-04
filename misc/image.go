package misc

import (
	"image"
	"io"
)

import (
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	_ "image/png"
)

type ImageInfo struct {
	Width  int
	Height int
}

func GetImageInfo(r io.Reader) (*ImageInfo, error) {
	img, _, err := image.DecodeConfig(r)
	if err != nil {
		return nil, err
	}

	return &ImageInfo{
		Width:  img.Width,
		Height: img.Height,
	}, nil
}
