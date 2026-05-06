package httplib

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	handlers []HandlerFunc
	index    int
	written  bool
}

type HandlerFunc func(c *Context)

// Next вызывает следующий обработчик в цепочке
// В middleware вызывайте Next() в конце, чтобы продолжить цепочку
func (c *Context) Next() {
	if c.index < len(c.handlers) {
		i := c.index
		c.index++
		c.handlers[i](c)
	}
}

// Param возвращает значение path-параметра по имени (например, "owner" для пути /repos/{owner}/{name})
func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

// JSON отправляет ответ с JSON и заданным статусом
func (c *Context) JSON(status int, v any) {
	if c.written {
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(c.Writer).Encode(v)
	}

	c.written = true
}

// String отправляет текстовый ответ
func (c *Context) String(status int, s string) {
	if c.written {
		return
	}

	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(s))
	c.written = true
}

// Error отправляет JSON {"error": "message"} с заданным статусом
func (c *Context) Error(message string, status int) {
	c.JSON(status, map[string]string{"error": message})
}

// Status выставляет HTTP-статус без тела. Полезно для 204 No Content
func (c *Context) Status(status int) {
	if c.written {
		return
	}

	c.Writer.WriteHeader(status)
	c.written = true
}

// Written возвращает true, если ответ уже отправлен
func (c *Context) Written() bool {
	return c.written
}

// Abort останавливает цепочку (дальнейшие обработчики не вызываются)
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// BindJSON декодирует тело запроса в v
// Возвращает ошибку при невалидном JSON
func (c *Context) BindJSON(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}
