package routes

import (
	"time"

	"github.com/catalinfl/gologin/database"
	"github.com/catalinfl/gologin/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Login(api *fiber.App) {
	route := api.Group("/login")

	route.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"test": "text"})
	})

	route.Post("/", func(c *fiber.Ctx) error {
		var creds Credentials

		err := c.BodyParser(&creds)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid body request"})
		}

		collection := database.GetClient().Database("gotest").Collection("users")
		filter := bson.M{"name": creds.Name}

		var user bson.M

		err = collection.FindOne(c.Context(), filter).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid user or password"})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "internal server error"})
			}
		}

		hashedPassword := user["password"].(string)

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user or password"})
		}

		token, exp, err := handlers.CreateJWTToken(&handlers.RequestUser{
			Name: user["name"].(string),
		})

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "Bearer " + token,
			Expires:  time.Unix(exp, 0),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{
			"message": "login successful",
		})

	})

}
