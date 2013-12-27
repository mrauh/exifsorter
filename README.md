# exifsorter

exifsorter reads images (jpg) from a directory, sorts them ascending by
their exif creation date and renames them according to a specified prefix.

[![Build Status](https://drone.io/github.com/mrauh/exifsorter/status.png)](https://drone.io/github.com/mrauh/exifsorter/latest)

## Installation

With [Go](http://www.golang.org) installed on your machine:

	$ go get github.com/mrauh/exifsorter

# Usage / Example

Just navigate to the directory where your jpg files are and run exifsorter. If 
you don't specify any flag, the current directory is used and you are asked to 
provide a prefix for the resulting filenames.
