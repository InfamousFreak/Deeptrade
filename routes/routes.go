package routes

import (
	"github.com/InfamousFreak/Deeptrade/backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {

	app.Post("/profile/create", controllers.CreateUserProfile)
	app.Get("/profile/:id/details", controllers.ShowUserProfile)
	app.Get("/profile/show", controllers.ShowProfiles)

	app.Get("/stock/:sym", controllers.GetStockData)
	app.Get("/news/:symbol", controllers.GetStockNews)

	app.Get("/sentiment/:symbol", controllers.GetSentimentFromNews)
	app.Get("analytics/:symbol", controllers.GetAnalytics)
	app.Get("backtest/:symbol", controllers.RunBacktest)
	app.Get("predict/:symbol", controllers.Prediction)

}