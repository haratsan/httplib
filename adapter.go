package httplib

import "net/http"

// Wrap превращает стандартный http.HandlerFunc в HandlerFunc для роутера
// Обработчик получает c.Writer и c.Request
func Wrap(h func(http.ResponseWriter, *http.Request)) HandlerFunc {
	return func(c *Context) {
		h(c.Writer, c.Request)
	}
}

// WrapHandler превращает http.Handler в HandlerFunc
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
