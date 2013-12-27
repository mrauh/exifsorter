package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// Image contains the file information and the exif creation date of the image.
type Image struct {
	FileInfo os.FileInfo
	ExifDate time.Time
}

// NewImage creates a new Image und reads the exif creation date of the image.
func NewImage(fi os.FileInfo) (*Image, error) {
	image := Image{FileInfo: fi}

	path := filepath.Join(dir, fi.Name())

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	x, err := exif.Decode(f)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	date, err := x.Get(exif.DateTimeOriginal)
	if err != nil {
		return nil, err
	}

	image.ExifDate, err = time.Parse("2006:01:02 15:04:05", date.StringVal())
	if err != nil {
		return nil, err
	}

	return &image, nil
}

// Images is a slice of Image and implements the sorter interface.
type Images []*Image

func (i Images) Len() int {
	return len(i)
}

func (i Images) Swap(x, y int) {
	i[x], i[y] = i[y], i[x]
}

func (i Images) Less(x, y int) bool {
	return i[x].ExifDate.Before(i[y].ExifDate)
}
