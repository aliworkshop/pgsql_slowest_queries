package application

import (
	"fmt"
	"github.com/aliworkshop/pgsql_slowest_queries/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Interface interface {
	RegisterRoutes(handler handler.HandlerFunc, path string, method Method)
	Start() error
	GetDB() *gorm.DB
}

type app struct {
	fiber  *fiber.App
	db     *gorm.DB
	config config
}

func NewApp() Interface {
	app := &app{fiber: fiber.New()}
	conf := InitConfig()
	err := conf.Unmarshal(&app.config)
	if err != nil {
		panic(err)
	}
	err = conf.Sub("sql").Unmarshal(&app.config.Sql)
	if err != nil {
		panic(err)
	}
	app.db = NewPostgreSqlConnection(app.config.Sql)

	return app
}

func (a *app) GetDB() *gorm.DB {
	return a.db
}

func (a *app) Start() error {
	return a.fiber.Listen(fmt.Sprintf("%s:%s", a.config.General.Address, a.config.General.Port))
}
