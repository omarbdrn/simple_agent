package rest

import (
	"encoding/json"
	"strconv"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/omarbdrn/simple-api/pkg/database"
	"github.com/omarbdrn/simple-api/pkg/server"
)

func RunRestServer(mqServer *server.Connection) {
	app := fiber.New()

	app.Get("/api/v1/shares/free", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Connecting to Message Queue
		db := database.GetDB()

		ip_ranges := []database.IPRange{}
		db.Model(&database.IPRange{}).Find(&ip_ranges, "taken = ?", true)

		for _, iprange := range ip_ranges {
			var question database.Question
			result := db.First(&question, "ip_range = ?", iprange.IPRange)
			if result.Error == nil {
				continue
			}

			question.IPRange = iprange.IPRange
			db.Create(&question)

			response := server.MQQuestion{
				QuestionID: fmt.Sprintf("%d", question.ID),
				Question:   iprange.IPRange,
			}

			jsonified_question, err := json.Marshal(response)
			if err != nil {
				continue // Failed to jsonify the question, ignore it.
			}

			mqServer.SendQuestion(mqServer.QuestionsChannel, server.Questions_channel, jsonified_question)
			go server.CheckQuestion(iprange.IPRange)
		}

		response := map[string]string{
			"status": "200",
		}

		return c.JSON(response)
	})

	app.Get("/api/v1/shares/get", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		maxCIDRSHeader := c.Get("Max-CIDRS")

		maxCIDRS, err := strconv.Atoi(maxCIDRSHeader)
		if err != nil {
			return c.Status(500).SendString("Invalid Max-CIDRS value")
		}

		db := database.GetDB()

		ip_ranges := []database.IPRange{}
		db.Model(&database.IPRange{}).Limit(maxCIDRS).Find(&ip_ranges, "taken = ?", false)

		if len(ip_ranges) == 0 {
			return c.Status(500).SendString("No IP Ranges found")
		}

		shareCIDRs := make([]string, len(ip_ranges))
		for i, ipRange := range ip_ranges {
			shareCIDRs[i] = ipRange.IPRange
		}

		share_model := &database.Share{IPRanges: ip_ranges}
		db.Create(share_model)

		// Marking IPRanges as Taken
		for _, ip_range := range ip_ranges {
			ip_range.Taken = true
			db.Save(&ip_range)
		}

		return c.JSON(&Share{ID: int(share_model.ID), Name: fmt.Sprintf("Share %d", share_model.ID), CIDRs: shareCIDRs})
	})

	app.Post("/api/v1/ipranges/add", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		ipRange := new(IPRangeDTO)
		if err := c.BodyParser(ipRange); err != nil {
			return c.SendStatus(500)
		}

		jsonResponse := map[string]string{
			"status": "200",
		}

		db := database.GetDB()
		db.Create(&database.IPRange{IPRange: ipRange.CIDR})

		return c.JSON(jsonResponse)
	})

	// app.Post("/api/v1/hosts")

	app.Listen("127.0.0.1:6453")
}
