package main

import (
	"bufio"
	"bytes"

	"github.com/ninetwentyfour/go-wkhtmltoimage"
)

func ProcessImage(bf *bufio.Reader) ([]byte, error) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(bf)

	c := wkhtmltoimage.ImageOptions{
		BinaryPath: "./wkhtmltoimage",
		Input:      "-",
		Format:     "png",
		Html:       buf.String(),
	}

	return wkhtmltoimage.GenerateImage(&c)
}
