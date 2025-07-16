package controllers 

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
)


type Predictions struct {
	Symbol string `json:"symbol"`
	Date string `json:"date"`
	Prediction string `json:"prediction"`
	Confidence float64 `json:"confidence"`
	Accuracy float64 `json:"accuracy"`
	F1Score float64 `json:"score"`
	FeaturesUsed []string `json:"features"`
	Error string `json:"error,omitempty"` //omitempty hides this field if empty
}

func Prediction(c *fiber.Ctx) error {

	symbol := c.Params("symbol")
	if symbol == "" {
		log.Println("Request failed: Missing query parameter.")

		return 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "'symbol' parameter is required.",
		})
	}

	//if the query is found, proceed with the python script
	cmd := exec.Command("python3", "ml/predict.py", symbol)

	log.Printf("Executing command: python3 ml/predict.py %s", symbol)
	output, err := cmd.CombinedOutput() //combined output gets both stdout and stderr
	if err != nil {
		errorMsg := fmt.Sprintf("Error executing python script: %s\nOutput: %s", err, string(output))
		log.Println(errorMsg)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errorMsg,
		})
	}

	var prediction Predictions  //unmarshalling the JSON output into our go model struct
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(output))), &prediction); err != nil {
		// this error will mean the script ran but the json output was not in the format or invalid JSON
		errorMsg := fmt.Sprintf("Error parsing JSON from Python script: %s\nRaw Output: %s", err, string(output))
		log.Println(errorMsg)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errorMsg,
		})
	}

	if prediction.Error != "" {
		errorMsg := fmt.Sprintf("Script returned an error: %s", prediction.Error)
		log.Println(errorMsg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorMsg,
		})
	}

	log.Printf("Successfully served prediction for %s", symbol)
	return c.Status(fiber.StatusOK).JSON(prediction)

}