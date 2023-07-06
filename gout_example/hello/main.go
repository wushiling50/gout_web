package main

import (
	"gout/gout"
	"net/http"
)

func main() {
	r := gout.Default()
	// r.Use(gout.Cors())
	r.Use(gout.Limiter(100))
	r.GET("/", func(c *gout.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	// r.GET("/html", func(c *gout.Context) {
	// 	c.HTML(http.StatusOK, `
	// 	<html><body><h2>Hello,world!</h2></body></html>`)
	// })
	r.POST("/a", func(c *gout.Context) {
		user := c.PostForm("user")
		c.JSON(http.StatusOK, gout.H{
			"user": user,
		})
	})
	// r.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *gout.Context) {
	// 	q1 := c.Param("param")
	// 	q2 := c.Param("param1")
	// 	c.String(http.StatusOK, "q1:%v,q2:%v", q1, q2)
	// })
	//
	v2 := r.Group("/b")
	v2.Use(gout.Cors())
	{
		v2.GET("/c/:name", func(c *gout.Context) {
			name := c.Param("name")
			c.String(http.StatusOK, "name:%s", name)
		})

		v2.GET("/c/user", func(c *gout.Context) {
			c.String(http.StatusOK, "name:%s", "user1")
		})
		v2.GET("/c/login", func(c *gout.Context) {
			c.String(http.StatusOK, "name:%s", "login1")
		})

		v3 := v2.Group("/g")
		v3.PUT("/d", func(c *gout.Context) {
			c.String(http.StatusOK, "put")
		})
		v3.DELETE("/e", func(c *gout.Context) {
			c.String(http.StatusOK, "delete")
		})
	}
	r.Run(":8080")
}
