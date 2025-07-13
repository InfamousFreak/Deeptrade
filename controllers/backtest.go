/*package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RunBackTest(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	cmd := exec.Command("python3", "ml/backtest.py", symbol)
	out, err != nil {
		fmt.Println("Python error:", err)
		fmt.Println("Output": string(out))
		return c.Status(500).JSON(fiber.Map{
			"error": "Pythons script failed",
			"details": string(out),
		})
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var result map[string]interface{}
	var debugMessages []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var jsonObj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonObj); err != nil {
			debugMessages = append(debugMessages, line)
			continue
		}

		if _, hasDebug := jsonObj["debug"]; hasDebug {
			debugMessages = append(debugMessages, jsonObj["debug"].(string))
		} else {
			result = jsonObj
		}
	}

	if err := scanner.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error reading Python output",
		})
	}

	if result == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "No valid backtest result found in Python output",
		})
	}

	//optionally include debug messages in the response
	in len(debugMessages) > 0 {
		result["debug_messages"] = debugMessages
	}

	return c.JSON(result)
}*/








package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RunBacktest(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	cmd := exec.Command("python3", "ml/backtest.py", symbol)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Python error:", err)
		fmt.Println("Output:", string(out))
		return c.Status(500).JSON(fiber.Map{
			"error":   "Python script failed",
			"details": string(out),
		})
	}

	// Parse each line as a separate JSON object
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var result map[string]interface{}
	var debugMessages []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var jsonObj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonObj); err != nil {
			// If it's not valid JSON, treat it as a debug message
			debugMessages = append(debugMessages, line)
			continue
		}

		// If it contains debug info, store it separately
		if _, hasDebug := jsonObj["debug"]; hasDebug {
			debugMessages = append(debugMessages, jsonObj["debug"].(string))
		} else {
			// This should be the main result
			result = jsonObj
		}
	}

	if err := scanner.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error reading Python output",
		})
	}

	if result == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "No valid backtest result found in Python output",
		})
	}

	// Optionally include debug messages in the response
	if len(debugMessages) > 0 {
		result["debug_messages"] = debugMessages
	}

	return c.JSON(result)
}