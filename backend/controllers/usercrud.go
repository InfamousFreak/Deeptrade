package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/InfamousFreak/Deeptrade/models"
	"github.com/InfamousFreak/Deeptrade/handlers"
	"github.com/InfamousFreak/Deeptrade/database"
	"errors"

)


func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	newUser := new(models.UserProfile)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "Error": "Username, email, and password are required",
        })
    }

	var existingUser models.UserProfile
    if err := database.Db.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "success": false,
            "message": "User with this email already exists",
        })
    }

	hashedPassword, err := handlers.HashPassword(newUser.Password)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to hash password"})
    }
    newUser.Password = hashedPassword

	createResult := database.Db.Create(&newUser)
    if createResult.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to create user",
            "error":   "Database error",
        })
    }

	claims := jwt.MapClaims{
        "ID":    newUser.ID, 
        "email": newUser.Email,
        "exp":   time.Now().Add(time.Hour * 12).Unix(),         
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(config.Secret))

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to generate token",
            "error":   "Internal server error",
        })
    }

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "User Created successfully",
        "data": fiber.Map{
            "user": sanitizeUserData(newUser),
            "token": t,
        },
    })
}

func sanitizeUserData(user *models.UserProfile) fiber.Map {
    return fiber.Map{
        "id":              user.ID,
        "name":            user.Name,
        "email":           user.Email,
        "country":         user.Country,
    }
}