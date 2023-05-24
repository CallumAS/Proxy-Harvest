package server

import (
	"github.com/gofiber/fiber/v2"
	"proxy-harvester/internal/pkg/checker"
	"proxy-harvester/internal/pkg/config"
	"proxy-harvester/internal/pkg/manager"
	"proxy-harvester/internal/pkg/model"
	"proxy-harvester/internal/pkg/scanner"
)

func Start() {
	app := fiber.New()

	app.Get("/results", func(c *fiber.Ctx) error {
		var exportType = c.Query("export")
		switch exportType {
		case "text":
			list := ""
			for s := range checker.Results {
				list += s + "\n"
			}
			return c.SendString(list)
		default:
			return c.JSON(checker.Results)
		}
	})
	app.Post("/settings", func(c *fiber.Ctx) error {
		payload := struct {
			Key      string         `json:"key"`
			Settings model.Settings `json:"settings"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if payload.Key != config.Config("API_KEY") {
			return c.JSON(map[string]string{"message": "go away"})
		}
		//set scanner settings
		manager.Scan = payload.Settings.Scanner.Active
		scanner.Ranges = payload.Settings.Scanner.Ranges
		scanner.Ports = payload.Settings.Scanner.Ports
		scanner.RangeTasks = payload.Settings.Scanner.Tasks

		//scraper Settings
		manager.Scrape = payload.Settings.Scraper.Active

		//checker
		checker.SetTasks(payload.Settings.Checker.Tasks)
		checker.Timeout = payload.Settings.Checker.Timeout

		//set payload info to managers
		return c.JSON(map[string]string{"message": "done!"})
	})

	app.Get("/settings", func(c *fiber.Ctx) error {
		var key = c.Query("key")
		if key != config.Config("API_KEY") {
			return c.JSON(map[string]string{"message": "go away"})
		}

		return c.JSON(model.Settings{
			Scanner: model.ScannerSettings{
				Ranges: scanner.Ranges,
				Ports:  scanner.Ports,
				Active: manager.Scan,
				Tasks:  scanner.RangeTasks,
			},
			Scraper: model.ScraperSettings{
				Active: manager.Scrape,
			},
			Checker: model.CheckerSettings{
				Tasks:   checker.Tasks,
				Timeout: checker.Timeout,
			},
		})
	})
	app.Listen(":3000")
}
