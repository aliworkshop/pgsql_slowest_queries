package application

import (
	"github.com/aliworkshop/pgsql_slowest_queries/handler"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

func (a *app) RegisterRoutes(handler handler.HandlerFunc, path string, method Method) {
	switch method {
	case GET:
		a.fiber.Get(path, handler)
	case POST:
		a.fiber.Post(path, handler)
	case PUT:
		a.fiber.Put(path, handler)
	case DELETE:
		a.fiber.Delete(path, handler)
	}
}
