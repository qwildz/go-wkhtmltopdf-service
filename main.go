package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var processorWg = &sync.WaitGroup{}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	processorWg.Add(1)
	go StartHttp(ctx, processorWg)

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until interrupted

	// // Handle shutdown
	cancelFunc()       // Signal cancellation to context.Context
	processorWg.Wait() // Block here until are workers are done

	fmt.Println("Bye!")
}

// func main() {
// 	html := "<html><head></head><body><p style='color:red;'>example</p></body></html>"
// 	c := wkhtmltoimage.ImageOptions{BinaryPath: "wkhtmltoimage", Input: "-", Format: "png", Html: html}
// 	log.Println("ad")
// 	_, err := wkhtmltoimage.GenerateImage(&c)

// 	log.Println(err)
// }
