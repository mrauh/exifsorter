// exifsorter reads images (jpg) from a directory, sorts them ascending by
// their exif creation date and renames them according to a specified prefix.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	dir string
	p   string
)

func init() {
	// Parse the flags.
	flag.StringVar(&dir, "dir", "", "Directory of the images to sort")
	flag.StringVar(&p, "p", "", "Prefix of the renamed files")
	flag.Parse()
}

func main() {
	// If dir is empty, use the current working directory.
	if dir == "" {
		dir, _ = os.Getwd()
	}

	// If prefix is empty, ask for it.
	if p == "" {
		fmt.Print("Enter prefix: ")
		_, err := fmt.Scanf("%s", &p)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	}

	// Read the source directory (the result is sorted alphanumerically)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory '%s': %v", dir, err)
	}

	images := Images{}

	// Get all jpeg files of dir and their exif creation date.
	for _, f := range files {
		ext := getExt(f.Name())

		if ext == ".jpg" || ext == ".jpeg" {
			image, err := NewImage(f)
			if err != nil {
				log.Fatal(err)
			}

			images = append(images, image)
		}
	}

	imgCount := len(images)
	if imgCount == 0 {
		log.Fatal("No (jpeg) images found.")
	}

	// Test for possible filename collisions
	for _, i := range images {
		checkFilenames(i.FileInfo.Name(), imgCount)
	}

	// Sort the images ascending by creation date.
	sort.Sort(images)

	// Rename the images according to the specified prefix.
	for idx, i := range images {
		ext := getExt(i.FileInfo.Name())
		fOld := filepath.Join(dir, i.FileInfo.Name())
		filename := createFilename(idx+1, ext)
		fNew := filepath.Join(dir, filename)

		err := os.Rename(fOld, fNew)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("The images have been renamed successfully!")
}

func getExt(name string) string {
	return strings.ToLower(filepath.Ext(name))
}

func checkFilenames(name string, imgCount int) {
	ext := getExt(name)

	for i := 0; i < imgCount; i++ {
		filename := createFilename(i+1, ext)

		if name == filename {
			log.Fatalf("Possible filename collision with '%s'. "+
				"Start again and use an other prefix.", name)
		}
	}
}

func createFilename(i int, ext string) string {
	return fmt.Sprintf("%s%d%s", p, i, ext)
}
