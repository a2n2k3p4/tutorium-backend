package main

import (
	"log"

	//module name "github.com/a2n2k3p4/tutorium-backend"
	"github.com/a2n2k3p4/tutorium-backend/config/dbserver" //store functions related to connecting to PostgreSQL
	"github.com/a2n2k3p4/tutorium-backend/handlers"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
)

// Before running the server, change config/dbserver/config.go to correct connection info

func main() {
	cfg := dbserver.NewConfig()

	db, err := dbserver.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	models.Migrate(db)

	app := fiber.New()

	handlers.AllRoutes(db, app) // Register admin routes
	// Define the /users route and handler inline

	// เริ่มจากการสร้าง user แล้วการใส่ learner , teacher , admin ให้ user
	// app.Post("/users/registerUser"       , func(c *fiber.Ctx) error {...}
	// app.Post("/users/:id/registerLearner", func(c *fiber.Ctx) error {...})
	// app.Post("/users/:id/registerTeacher", func(c *fiber.Ctx) error {...})
	// app.Post("/users/:id/registerAdmin"  , func(c *fiber.Ctx) error {...})

	log.Fatal(app.Listen(":8000")) // using PORT 8000 (localhost:8000)
}
