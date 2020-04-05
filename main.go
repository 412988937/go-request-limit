package main

import (
	"github.com/gin-gonic/gin"
	"github.com/412988937/go-request-limit/limit"
	"log"
)

func main(){
	r := gin.Default()
	limiter, err := limit.NewLimiter("redis://localhost:6379", 1000, 60 * 60)
	if err != nil {
		log.Fatal(err)
	}
	r.Use(limit.RequestLimitMiddleware(limiter))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}