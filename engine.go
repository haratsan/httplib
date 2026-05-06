package httplib

import (
	"net/http"
	"strings"
)

type Engine struct {
	mux            *http.ServeMux
	globalHandlers []HandlerFunc
}

func New() *Engine {
	return &Engine{
		mux:            http.NewServeMux(),
		globalHandlers: nil,
	}
}

// Use добавляет глобальные middleware (вызываются перед каждым обработчиком маршрута)
func (e *Engine) Use(handlers ...HandlerFunc) {
	e.globalHandlers = append(e.globalHandlers, handlers...)
}

// Group возвращает группу с префиксом path. Middleware группы выполняются после глобальных
func (e *Engine) Group(path string) *Group {
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = ""
	}

	return &Group{
		engine:     e,
		prefix:     path,
		middleware: nil,
	}
}

// handle регистрирует маршрут method + path с цепочкой обработчиков (последний - финальный handler)
func (e *Engine) handle(method, path string, handlers []HandlerFunc) {
	if len(handlers) == 0 {
		panic("httplib: требуется как минимум один обработчик")
	}

	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	}

	pattern := method + " " + path
	chain := make([]HandlerFunc, 0, len(e.globalHandlers)+len(handlers))
	chain = append(chain, e.globalHandlers...)
	chain = append(chain, handlers...)
	e.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			Request:  r,
			Writer:   w,
			handlers: chain,
			index:    0,
		}
		c.Next()
	})
}

// GET регистрирует обработчик для GET
func (e *Engine) GET(path string, handlers ...HandlerFunc) {
	e.handle(http.MethodGet, path, handlers)
}

// POST регистрирует обработчик для POST
func (e *Engine) POST(path string, handlers ...HandlerFunc) {
	e.handle(http.MethodPost, path, handlers)
}

// PATCH регистрирует обработчик для PATCH.
func (e *Engine) PATCH(path string, handlers ...HandlerFunc) {
	e.handle(http.MethodPatch, path, handlers)
}

// PUT регистрирует обработчик для PUT
func (e *Engine) PUT(path string, handlers ...HandlerFunc) {
	e.handle(http.MethodPut, path, handlers)
}

// DELETE регистрирует обработчик для DELETE
func (e *Engine) DELETE(path string, handlers ...HandlerFunc) {
	e.handle(http.MethodDelete, path, handlers)
}

// ServeHTTP реализует http.Handler
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mux.ServeHTTP(w, r)
}

// Run запускает сервер на addr (удобная обёртка над http.ListenAndServe)
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
