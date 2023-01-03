package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

var singlePageHtml = `
<!doctype html>
<html lang="en">

<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <title>%s</title>

    <style>
        html,
        body {
            text-align: center;
            margin: 0;
        }

        img {
            page-break-inside: avoid;
        }

        img {
            max-width: 100%% !important;
            max-height: 345mm;
        }
    </style>
</head>
<body class="nobreak">
    <img class="nobreak auto-page-break" src="data:image/png;base64, %s">
</body>
</html>`

func ProcessPdf(title string, bf *bufio.Reader) ([]byte, error) {

	// Create new PDF generator
	wkhtmltopdf.SetPath("./wkhtmltopdf")
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

func ProcessPdfSinglePage(title string, bf *bufio.Reader) ([]byte, error) {

	// Convert into image first
	data, err := ProcessImage(bf)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can't convert into image")
	}

	generatedHtml := fmt.Sprintf(singlePageHtml, title, base64.RawStdEncoding.EncodeToString(data))
	reader := bufio.NewReader(bytes.NewBufferString(generatedHtml))

	return ProcessPdf(title, reader)
}
