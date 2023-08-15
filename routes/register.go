package routes

import (
	"time"

	"github.com/catalinfl/gologin/database"
	"github.com/catalinfl/gologin/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(app *fiber.App) {

	route := app.Group("/register")

	route.Post("/", func(c *fiber.Ctx) error {
		req := new(handlers.RequestUser)

		if err := c.BodyParser(req); err != nil {
			return err
		}

		if req.Name == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid signup credentials")
		}

		collection := database.GetClient().Database("gotest").Collection("users")
		filter := bson.M{"name": req.Name}
		name := collection.FindOne(c.Context(), filter)

		if name.Err() == nil {
			return fiber.NewError(fiber.StatusNotAcceptable, "User already exists")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		user := &handlers.RequestUser{
			Name:     req.Name,
			Password: string(hash),
		}

		collection = database.GetClient().Database("gotest").Collection("users")

		_, err = collection.InsertOne(c.Context(), user)

		if err != nil {
			return err
		}

		token, exp, err := handlers.CreateJWTToken(user)

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "Bearer " + token,
			Expires:  time.Unix(exp, 0),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	route.Get("/protected", handlers.VerifyToken, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"test": "test1"})
	})
}
