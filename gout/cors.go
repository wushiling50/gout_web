package gout

import (
	"net/http"
)

func Cors() HandlerFunc {
	return func(c *Context) {
		method := c.Req.Method               //请求方法
		origin := c.Req.Header.Get("Origin") //请求头部
		if origin != "" {
			// 这是允许访问所有域
			c.SetHeader("Access-Control-Allow-Origin", "*")
			//  跨域请求是否需要带cookie信息 默认设置为true
			c.SetHeader("Access-Control-Allow-Credentials", "true")
			//获取其他字段
			c.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// 表明服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段。
			c.SetHeader("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 缓存请求信息 单位为秒
			c.SetHeader("Access-Control-Max-Age", "172800")

		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next() //  处理请求
	}
}
