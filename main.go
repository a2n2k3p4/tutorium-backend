package main

import (
	"log"

	// module name "github.com/a2n2k3p4/tutorium-backend"
	"github.com/a2n2k3p4/tutorium-backend/config/dbserver" // store functions related to connecting to PostgreSQL
	"github.com/a2n2k3p4/tutorium-backend/handlers"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// swagger
	_ "github.com/a2n2k3p4/tutorium-backend/docs"
	"github.com/gofiber/swagger"
)

// Before running the server, change config/dbserver/config.go to correct connection info

//	@title			Tutorium Backend API
//	@version		1.0
//	@description	This is the API for Tutorium Backend system.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	AGPL-3.0
//	@license.url	https://www.gnu.org/licenses/agpl-3.0.en.html

func main() {
	cfg := dbserver.NewConfig()

	db, err := dbserver.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	models.Migrate(db)

	// path
	app := fiber.New()

	app.Use(cors.New())

	// debug route
	app.Get("/", func(c *fiber.Ctx) error {
		log.Printf("c base url : %s", c.BaseURL())
		return c.JSON(fiber.Map{
			"message": "Tutorium Backend API",
			"swagger": c.BaseURL() + "/swagger/",
		})
	})

	// custom swagger UI
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:                      "doc.json", // swagger.json location
		DeepLinking:              true,
		DocExpansion:             "false",
		DefaultModelsExpandDepth: 2, // expand models
	}))

	handlers.AllRoutes(app) // Register admin routes
	// Define the /users route and handler inline

	// เริ่มจากการสร้าง user แล้วการใส่ learner , teacher , admin ให้ user
	// app.Post("/users/registerUser"       , func(c *fiber.Ctx) error {...}
	// app.Post("/users/:id/registerLearner", func(c *fiber.Ctx) error {...})
	// app.Post("/users/:id/registerTeacher", func(c *fiber.Ctx) error {...})
	// app.Post("/users/:id/registerAdmin"  , func(c *fiber.Ctx) error {...})

	// lg
	log.Println("Server starting on :8000")
	log.Println("API endpoint: /")
	log.Println("Swagger UI: /swagger/")
	log.Fatal(app.Listen(":8000"))
}
