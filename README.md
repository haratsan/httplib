# httplib

Минимальная Go-библиотека для HTTP-роутинга и промежуточных обработчиков

```go
e := httplib.New()
e.Use(Logger, Recovery)

api := e.Group("/api")
api.POST("/register", Wrap(registerHandler))
api.GET("/user", RequireAuth(jwt), Wrap(currentUserHandler))

admin := api.Group("/admin")
admin.Use(RequireAdmin)
admin.GET("/users", Wrap(GetUsers))

http.ListenAndServe(":8000", e)
```