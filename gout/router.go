package gout

import (
	"net/http"
	"strings"
)

type router struct {
	managers map[string]*Manager    //设置分组，组名为method   (EX:managers['GET'],managers['POST'])
	handlers map[string]HandlerFunc //handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
}

func newRouter() *router {
	return &router{
		managers: make(map[string]*Manager),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将pattern进行解析
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
		}
	}

	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	//查找manager
	_, ok := r.managers[method]
	if !ok {
		e := make([]Entry, 0)
		r.managers[method] = &Manager{
			menbers: &e,
		}
	}
	r.managers[method].insert(pattern, parts)
	r.handlers[key] = handler
}

// 查找路由
func (r *router) getRoute(method string, path string) (string, map[string]string) {
	parts := parsePattern(path)
	params := make(map[string]string)
	manager, ok := r.managers[method]

	if !ok {
		return "", nil
	}
	path, paramsmap := manager.search(path, parts)
	if path != "" && paramsmap != nil {
		for key, part := range paramsmap {
			if part == "" {
				continue
			}
			if part[0] == ':' {
				params[part[1:]] = parts[key]
			}
		}
	}

	return path, params
}

func (r *router) handle(c *Context) {
	path, params := r.getRoute(c.Method, c.Path)
	if path != "" {
		key := c.Method + "-" + path
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
