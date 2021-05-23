package main

import (
	"bufio"
	"errors"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func ProcessHtml(title string, bf *bufio.Reader) ([]byte, error) {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can't initialize wkhtmltopdf")
	}

	pdfg.Title.Set(title)

	// Add to document
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bf))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can't create pdf")
	}

	log.Println("Done")

	return pdfg.Bytes(), nil
}
