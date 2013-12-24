package main

import (
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	dir string
	p   string
)

//--------------------------------------------------------------------------

func init() {
	flag.StringVar(&dir, "dir", "", "Source directory of your images")
	flag.StringVar(&p, "p", "image", "Prefix of the renamed files")
	flag.Parse()
}

//--------------------------------------------------------------------------

// Image contains the file information and the date of the .
type Image struct {
	FileInfo os.FileInfo
	Date     time.Time
}

// Images is a slice of Image and implements the sorter interface.
type Images []Image

func (i Images) Len() int {
	return len(i)
}

func (i Images) Swap(x, y int) {
	i[x], i[y] = i[y], i[x]
}

func (i Images) Less(x, y int) bool {
	return i[x].Date.Before(i[y].Date)
}

//--------------------------------------------------------------------------

func main() {
	// If dir is empty, use the current working directory.
	if dir == "" {
		dir, _ = os.Getwd()
	}

	// If prefix is ommited, use 'image'.
	if p == "" {
		p = "image"
	}

	// Read the source directory (the result is sorted alphanumeric)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory '%v': %v", dir, err)
	}

	images := Images{}

	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))

		if ext == ".jpg" || ext == ".jpeg" {
			image := Image{FileInfo: f}

			image.Date, err = getDate(image.FileInfo)
			if err != nil {
				log.Fatal(err)
			}

			images = append(images, image)
		}
	}

	// Sort the images ascending by creation date.
	sort.Sort(images)

	// Rename the images according to the specified format.
	for idx, i := range images {
		ext := strings.ToLower(filepath.Ext(i.FileInfo.Name()))

		fOld := filepath.Join(dir, i.FileInfo.Name())

		filename := fmt.Sprintf("scheidegg%v%v", idx+1, ext)
		fNew := filepath.Join(dir, filename)

		err := os.Rename(fOld, fNew)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("ok")
}

//--------------------------------------------------------------------------

func getDate(file os.FileInfo) (time.Time, error) {
	var t time.Time

	path := filepath.Join(dir, file.Name())

	f, err := os.Open(path)
	if err != nil {
		return t, err
	}

	x, err := exif.Decode(f)
	defer f.Close()
	if err != nil {
		return t, err
	}

	date, err := x.Get(exif.DateTimeOriginal)
	if err != nil {
		return t, err
	}

	t, err = time.Parse("2006:01:02 15:04:05", date.StringVal())
	if err != nil {
		return t, err
	}

	return t, nil
}

//--------------------------------------------------------------------------
