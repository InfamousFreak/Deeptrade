package repository

import (
	"errors"
	"github.com/InfamousFreak/Deeptrade/models"
	"gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

/*func FindbyCredentials(email, password string) (*models.UserProfile, error) {

	if email == "test@gmail.com" && password == "test12345" {
		return &models.UserProfile{
			Email: "test@gmail.com",
			Password: "test12345",
			Country: "India",
		}, nil
	}
	return nil, errors.New("user not found")
}*/

func Find(db *gorm.DB, email, password string) (*models.UserProfile, error) {
    var user models.UserProfile

    // Query the database for the user with the given email
    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("Invalid Credentials")
        }
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("Invalid Credentials")
    }

    return &user, nil
}
