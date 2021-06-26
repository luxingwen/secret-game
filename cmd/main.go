package main

import (
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/luxingwen/secret-game/controller"
	_ "github.com/luxingwen/secret-game/dao"
)

func main() {
	fmt.Println("hello world")

	gin.SetMode(gin.DebugMode)
	go controller.WsManager.Start()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        controller.Router(),
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("start server  on ", 8080)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
