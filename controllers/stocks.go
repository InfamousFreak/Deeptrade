package controllers

import (
		"github.com/gofiber/fiber/v2"
		"os/exec"
		"encoding/json"
		"os"
		"fmt"
		"net/http"
		"time"
		"math/rand"
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

func ChartHandler(c *fiber.Ctx) error {
	// In a real application, you would fetch real historical data here.
	// This is just generating random sample data for demonstration.
	var seriesData []fiber.Map
	currentTime := time.Now()
	currentPrice := 150.0 + rand.Float64()*20

	for i := 0; i < 60; i++ {
		date := currentTime.AddDate(0, 0, -i)
		seriesData = append(seriesData, fiber.Map{
			"x": date.Format("2006-01-02"),
			"y": []float64{
				currentPrice - rand.Float64()*2, // Open
				currentPrice + rand.Float64()*2, // High
				currentPrice - rand.Float64()*3, // Low
				currentPrice + rand.Float64(),   // Close
			},
		})
		currentPrice += (rand.Float64() - 0.5) * 5 // Fluctuate price
	}

	// Reverse the data so it's in chronological order
	for i, j := 0, len(seriesData)-1; i < j; i, j = i+1, j-1 {
		seriesData[i], seriesData[j] = seriesData[j], seriesData[i]
	}

	return c.JSON(fiber.Map{
		"series": []fiber.Map{
			{"data": seriesData},
		},
	})
}


