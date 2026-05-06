package httplib

import (
	"net/http"
	"strings"
)

type Group struct {
	engine     *Engine
	prefix     string
	middleware []HandlerFunc
}

// Use добавляет middleware к группе (выполняются после глобальных, перед обработчиком маршрута)
func (g *Group) Use(handlers ...HandlerFunc) {
	g.middleware = append(g.middleware, handlers...)
}

// Group возвращает вложенную группу с префиксом path (относительно текущей группы)
func (g *Group) Group(path string) *Group {
	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	newPrefix := g.path(path)
	if newPrefix == "" {
		newPrefix = "/"
	}

	return &Group{
		engine:     g.engine,
		prefix:     newPrefix,
		middleware: nil,
	}
}

func (g *Group) path(p string) string {
	p = strings.TrimPrefix(p, "/")
	if g.prefix == "" || g.prefix == "/" {
		return "/" + p
	}

	return strings.TrimSuffix(g.prefix, "/") + "/" + p
}

func (g *Group) handle(method, path string, handlers []HandlerFunc) {
	if len(handlers) == 0 {
		panic("httplib: требуется как минимум один обработчик")
	}

	fullPath := g.path(path)
	chain := make([]HandlerFunc, 0, len(g.middleware)+len(handlers))
	chain = append(chain, g.middleware...)
	chain = append(chain, handlers...)
	g.engine.handle(method, fullPath, chain)
}

// GET регистрирует GET в группе
func (g *Group) GET(path string, handlers ...HandlerFunc) {
	g.handle(http.MethodGet, path, handlers)
}

// POST регистрирует POST в группе
func (g *Group) POST(path string, handlers ...HandlerFunc) {
	g.handle(http.MethodPost, path, handlers)
}

// PATCH регистрирует PATCH в группе
func (g *Group) PATCH(path string, handlers ...HandlerFunc) {
	g.handle(http.MethodPatch, path, handlers)
}

// PUT регистрирует PUT в группе
func (g *Group) PUT(path string, handlers ...HandlerFunc) {
	g.handle(http.MethodPut, path, handlers)
}

// DELETE регистрирует DELETE в группе
func (g *Group) DELETE(path string, handlers ...HandlerFunc) {
	g.handle(http.MethodDelete, path, handlers)
}
