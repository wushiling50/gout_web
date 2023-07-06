package main

import (
	"gout/gout"
	"net/http"
	"time"
)

func main() {
	r := gout.Default()

	r.GET("/", func(c *gout.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	r.GET("/1", func(c *gout.Context) {
		c.String(http.StatusOK, "Hello World")
		time.Sleep(time.Second * time.Duration(10))
	})
	r.Run(":8080")
}
