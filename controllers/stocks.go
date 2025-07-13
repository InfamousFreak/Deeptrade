package controllers

import (
		"github.com/gofiber/fiber/v2"
		"os/exec"
		"encoding/json"
		"os"
		"fmt"
		"net/http"
)

func GetStockData(c *fiber.Ctx) error {
	symbol := c.Params("sym")

	out, err := exec.Command("python3", "ml/fetch_stock.py", symbol).Output()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch stock data",
		})
	}

	var data []map[string]interface{}
	json.Unmarshal(out, &data)

	return c.JSON(data)
}

func GetStockNews(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	apiKey := os.Getenv("NEWS_API_KEY")

	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&sortBy=publishedAt&language=en&apiKey=%s", symbol, apiKey)

	res, err := http.Get(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "News fetch failed"})
	}
	defer res.Body.Close()

	var news map[string]interface{}
	json.NewDecoder(res.Body).Decode(&news)

	return c.JSON(news)
}


