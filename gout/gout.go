package gout

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var CloseCh = make(chan struct{}, 100)

// 提供给框架用户的，用来定义路由映射的处理方法,接受Context
type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // 支持中间件
		parent      *RouterGroup  // 支持嵌套
		engine      *Engine       // 所有分组共享一个Engine实例
	}

	// 结构体Engine，作为Handler接口的实现类进行注入
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // 储存所有分组
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// Group被定义为创建一个新的RouterGroup，所有组共享同一个Engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 嵌套addRoute方法调用，添加路由进表
func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 添加GET请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 添加POST请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 添加PUT请求
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}

// 添加DELETE请求
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

// 启动一个Http服务
func (engine *Engine) Run(addr string) {
	s := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	go s.ListenAndServe()
	graceful_shutdown(s)
}

// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) { //通过路由前缀判断引用哪些中间件
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares

	engine.router.handle(c)
}

func graceful_shutdown(srv *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case signal := <-sig:
			switch signal {
			case syscall.SIGHUP, syscall.SIGINT:
				log.Printf("Received signal:%s\n", signal)
				CloseCh <- struct{}{}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					log.Fatal("Server Shutdown:", err)
				}
				log.Println("exit!")
				return
			}
		default:
			time.Sleep(time.Microsecond * time.Duration(1))
		}
	}

}
