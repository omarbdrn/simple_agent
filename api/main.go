package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Share struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	CIDRs []string `json:"cidrs"`
}

func main() {
	app := fiber.New()
	cidrs := []string{}

	app.Get("/api/v1/shares/get", func(c *fiber.Ctx) error {
		maxCIDRsHeader := c.Get("Max-CIDRs")
		maxCIDRs, err := strconv.Atoi(maxCIDRsHeader)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Max-CIDRs header")
		}

		if maxCIDRs > len(cidrs) {
			maxCIDRs = len(cidrs)
		}

		share := Share{
			ID:    1,
			Name:  "Share 1",
			CIDRs: cidrs[:maxCIDRs],
		}

		return c.JSON(share)
	})

	app.Post("/api/v1/hosts", func(c *fiber.Ctx) error {
		var requestBody map[string]interface{}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON format")
		}

		log.Printf("Received JSON: %+v", requestBody)

		return c.JSON(map[string]string{"ping": "pong"})
	})

	log.Fatal(app.Listen("127.0.0.1:6453"))
}
