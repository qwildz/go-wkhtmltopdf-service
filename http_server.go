package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func StartHttp(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	router := gin.Default()
	registerRoute(router)

	srv := &http.Server{
		Addr:    ":8787",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down http server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

func registerRoute(router *gin.Engine) {
	router.GET("/", routerStatus)
	router.POST("/convert", routerProcessHtml)
}

func routerStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func routerProcessHtml(c *gin.Context) {
	file, err := c.FormFile("html")
	if err != nil {
		log.Println("Can't process request. No file posted.")
	}

	f, err := file.Open()
	bf := bufio.NewReader(f)

	defer f.Close()

	data, err := ProcessHtml(c.PostForm("title"), bf)
	if err != nil {
		log.Println("Can't process request.")
		c.JSON(500, gin.H{
			"message":        "Can't process request",
			"message_detail": err,
		})
		return
	}

	c.Data(200, "application/pdf", data)
}
