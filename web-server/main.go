package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func createApp() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "world")
	})
	return r
}

func main() {
	a := func() string {
		if a, ok := os.LookupEnv("LISTEN_ADDRESS"); ok {
			return a
		}
		return "0.0.0.0:8080"
	}()
	s := createApp()
	err := s.Run(a)
	fmt.Fprintf(os.Stderr, "failed to start server: %s", err)
}
