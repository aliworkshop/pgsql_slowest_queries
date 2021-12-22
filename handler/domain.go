package handler

import "github.com/gofiber/fiber/v2"

type HandlerFunc func(c *fiber.Ctx) error

type Interface interface {
	Hello(c *fiber.Ctx) error
	SlowestConnectionHandler(c *fiber.Ctx) error
}
