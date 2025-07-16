package controllers

import ( 
	"encoding/json"
	"fmt"
	"os/exec"
	"os"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/*var sampleHeadlines = []string{
	"Apple stock hits record high after earnings beat expectations",
	"Market shows mixed reaction to inflation report",
	"Tesla faces production delays in Shanghai plant",
}*/

func GetSentimentFromNews(c *fiber.Ctx) error {

	symbol := c.Params("symbol")
	apiKey := os.Getenv("NEWS_API_KEY")
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&sortBy=publishedAt&language=en&apiKey=%s", symbol, apiKey)

	// Fetch news
	resp, err := http.Get(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch news"})
	}
	defer resp.Body.Close()

	var news map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&news)

	//extract healines
	articles, ok := news["articles"].([]interface{})
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid article format"})
	}

	var headlines []string
	for _, a := range articles {
		article := a.(map[string]interface{})
		title, ok := article["title"].(string)
		if ok {
			headlines = append(headlines, title)
		}
	}

	//convert headlines to json
	headlineBytes, _ := json.Marshal(headlines)

	//run python scrit for headliens
	cmd := exec.Command("python3", "ml/sentiment.py", string(headlineBytes))
	out, err := cmd.CombinedOutput()
	if err != nil {
    	fmt.Println("❌ Python error:", err)
    	fmt.Println("⚠️ Python stderr:", string(out)) // shows actual error
    	return c.Status(500).JSON(fiber.Map{"error": "Python script failed", "details": string(out)})
}


	//unmarshal result from python
	var result []map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid JSON from sentiment engine"})
	}


	return c.JSON(fiber.Map{
		"success": true,
		"symbol": symbol,
		"analyzed": result,
	})
}

