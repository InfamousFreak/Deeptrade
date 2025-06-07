package handlers

import (
	"time"
	"github.com/gofiber/fiber/v2"
 jtoken "github.com/golang-jwt/jwt/v4"
	"github.com/InfamousFreak/Deeptrade/config"
	"github.com/InfamousFreak/Deeptrade/models"
	"github.com/InfamousFreak/Deeptrade/repository"
	"github.com/InfamousFreak/Deeptrade/database"
	
)

//login rouet
func Login(c *fiber.Ctx) error {
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
  		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
   			"error": err.Error(),
  	})
 }

 user, err := repository.Find(database.Db, loginRequest.Email, loginRequest.Password)
 if err != nil {
  return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
   "error": err.Error(),
  })
 }

 day := time.Hour * 24
 // Create the JWT claims, which includes the user ID and expiry time
 claims := jtoken.MapClaims{
  "ID":		   user.ID,
  "email": 	   user.Email,
  "country":   user.Country,
  "exp":   time.Now().Add(day * 1).Unix(),
 }

 token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
 // Generate encoded token and send it as response.
 t, err := token.SignedString([]byte(config.Secret))
 if err != nil {
  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
   "error": err.Error(),
  })
 }

 return c.JSON(models.LoginResponse{
  Token: t,
 })
}

func Protected(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	email := claims["email"].(string)
	country := claims["country"].(string)
	return c.SendString("Welcome to protected route, access provided. Your email is " + email + "and you are from " + country)
}