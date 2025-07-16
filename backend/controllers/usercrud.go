package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/InfamousFreak/Deeptrade/backend/models"
	"github.com/InfamousFreak/Deeptrade/backend/database"
    "github.com/InfamousFreak/Deeptrade/backend/passwordhashing"
    jwt "github.com/golang-jwt/jwt/v4"
	"errors"
    "time"
    "github.com/InfamousFreak/Deeptrade/backend/config"

)


/*func CreateUser(c *fiber.Ctx) error {
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

	hashedPassword, err := passwordhashing.HashPassword(newUser.Password)
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
}*/

func CreateUserProfile(c *fiber.Ctx) error {
	newUser := new(models.UserProfile)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	// Validate required fields
	if newUser.Email == "" || newUser.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "Email and password are required",
		})
	}

	// Check if user already exists
	var existingUser models.UserProfile
	if err := database.Db.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "User with this email already exists",
		})
	}

	// Hash password
	hashedPassword, err := passwordhashing.HashPassword(newUser.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to hash password"})
	}
	newUser.Password = hashedPassword

	// Create user in database
	createResult := database.Db.Create(&newUser)
	if createResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   "Database error",
		})
	}

	// Generate JWT token
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

	// Return success response with user data and token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"data": fiber.Map{
			"user":  sanitizeUserData(newUser),
			"token": t,
		},
	})
}

func sanitizeUserData(user *models.UserProfile) fiber.Map {
	return fiber.Map{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"country": user.Country,
	}
}

func ShowUserProfile(c *fiber.Ctx) error {
    userID := c.Params("id")
    if userID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "User ID is required",
        })
    }

    var user models.UserProfile
    result := database.Db.First(&user, userID)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "User not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error fetching user profile",
        })
    }

    // Sanitize user data before sending
    sanitizedUser := sanitizeUserData(&user)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "User profile fetched successfully",
        "data":    sanitizedUser,
    })
}

func ShowProfiles(c *fiber.Ctx) error {
    var users []models.UserProfile

    // Query the database for all user profiles
    result := database.Db.Find(&users)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(&fiber.Map{
        "data":    users,
        "success": true,
        "message": "Retrieved Successfully",
    })
}