package config

import (
	"fmt"
	"os"

	 "github.com/joho/godotenv"
)

/*func Config(key string) string {
    // load .env file
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Print("Error loading .env file")
    }
        
    return os.Getenv(key)
}*/

func Load(key string) string {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	//return the value of the variable
	return os.Getenv(key)
}

// secret key used to sign the JWT, this must be a secure key and should not be stored in the code
const Secret = "secret"


