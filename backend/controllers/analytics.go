package controllers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"github.com/gofiber/fiber/v2"
)

func GetAnalytics(c *fiber.Ctx) error {
	symbol := c.Params("symbol")

	if symbol == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Stock symbol is required"})
	}


	cmd := exec.Command("python3", "ml/analytics.py", symbol)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("❌ Python error:", err)
		fmt.Println("⚠️ Output:", string(out))
		return c.Status(500).JSON(fiber.Map{
			"error": "Python script failed",
			"details": string(out),
		})
	}

	cleanOutput := strings.TrimSpace(string(out))

	start := strings.Index(cleanOutput, "{")
	end := strings.LastIndex(cleanOutput, "}")
	
	if start == -1 || end == -1 || start >= end {
		return c.Status(500).JSON(fiber.Map{
			"error": "No valid JSON found in output",
			"raw_output": cleanOutput,
		})
	}
	
	jsonStr := cleanOutput[start:end+1]

	fmt.Printf("🔍 Extracted JSON: %s\n", jsonStr)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Println("JSON Parse error:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid JSON from analytics engine",
			"raw_output": cleanOutput,
			"extracted_json": jsonStr,
			"parse_error": err.Error(),
		})
	}

	return c.JSON(result)
}