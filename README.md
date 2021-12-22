# pgsql_slowest_queries
postgresql slowest connected database queries.
this repo uses fiber as a go framework and postgresql as a db

### First
get project 
```shell script
git clone https://github.com/aliworkshop/pgsql_slowest_queries.git
```

### Load Modules
```shell script
go get -v ./...
```

### Run
first update config.yaml with correct information then run this command
```shell script
go run main.go
```

### Test
```shell script
go test ./...
```

### Develop
for develop just define handler with HandlerFunc type and 
add method to handler interface and register route for that handler in main.go

```shell script
package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *handler) Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
```

```shell script
package main

import (
	"github.com/aliworkshop/pgsql_slowest_queries/application"
	"github.com/aliworkshop/pgsql_slowest_queries/handler"
	"log"
)

func main() {

	app := application.NewApp()
	handlers := handler.NewHandler(app.GetDB())

	app.RegisterRoutes(handlers.Hello, "/", application.GET)

	if err := app.Start(); err != nil {
		log.Fatal("error happened on app starting... : " + err.Error())
	}
}

```
