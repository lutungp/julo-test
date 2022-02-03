package accountHandler

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lutungp/julo-test/database"
	"github.com/lutungp/julo-test/internal/handler/utils"
	"github.com/lutungp/julo-test/internal/model"
)

type ResponseToken struct {
	Status string    `json:"status"`
	Data   DataToken `json:"data"`
}

type DataToken struct {
	Token string `json:"token"`
}

type SomeStruct struct {
	Name string
	Age  uint8
}

func CreateAccount(c *fiber.Ctx) error {
	db := database.DB
	var account model.Account

	err := c.BodyParser(account)

	id := c.Params("customerXid")

	db.Find(&account, "customer_xid = ?", id)

	if account.ID == uuid.Nil {
		account.ID = uuid.New()
		err = db.Create(&account).Error

		if err == nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create account", "data": err})
		}
	}

	claims := jwt.MapClaims{
		"ID":          account.ID,
		"customerXid": account.CustomerXid,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(token)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return c.JSON(fiber.Map{"status": "failure", "message": "Can't generate token"})
	}

	res := &ResponseToken{
		Status: "success",
		Data: DataToken{
			Token: t,
		},
	}

	return c.JSON(res)
}

func GetNewAccessToken(c *fiber.Ctx) error {
	// Generate a new Access token.
	token, err := utils.GenerateNewAccessToken()
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":        false,
		"msg":          nil,
		"access_token": token,
	})
}
