package gout

import (
	"net/http"
	"testing"
)

func BenchmarkGetRoute(B *testing.B) {
	r := New()
	r.GET("/", func(c *Context) {})
	runRequest(B, r, "GET", "/")
}

func BenchmarkDefault(B *testing.B) {
	r := Default()
	r.GET("/", func(c *Context) {})
	runRequest(B, r, "GET", "/")
}

func BenchmarkCorsMiddleware(B *testing.B) {
	r := New()
	r.Use(Cors())
	r.GET("/", func(c *Context) {})
	runRequest(B, r, "GET", "/")
}

func BenchmarkFiveRoute(B *testing.B) {
	r := New()
	v1 := r.Group("/ping")
	v1.Use(Logger(), Recovery())
	{
		r.GET("/get", func(ctx *Context) {})
		r.POST("/post", func(ctx *Context) {})
		r.PUT("/put", func(ctx *Context) {})
		r.DELETE("delete", func(ctx *Context) {})
	}
}

func BenchmarkParmes(B *testing.B) {
	r := New()
	r.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {
		q1 := c.Param("param")
		q2 := c.Param("param1")
		c.String(http.StatusOK, "q1:%v,q2:%v", q1, q2)
	})
	runRequest(B, r, "GET", "/param/path/to/parameter/john/12345")
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}
func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}
func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func (m *mockWriter) WriteHeader(int) {}

func runRequest(B *testing.B, r *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}
