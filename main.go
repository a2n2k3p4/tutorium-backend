package main

import (
	"context"
	"log"

	//module name "github.com/Parkorn/KUTutorium"
	"github.com/Parkorn/KUTutorium/config/dbserver" //store functions related to connecting to PostgreSQL
	"github.com/gofiber/fiber/v2"
)

//Before running the server, change config/dbserver/config.go to correct connection info

func main() {
	cfg := dbserver.NewConfig()

	db, err := dbserver.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	// Define the /users route and handler inline
	app.Get("/users", func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(), "SELECT user_id, session_id, balance FROM users")
		if err != nil {
			return c.Status(500).SendString("Failed to query users")
		}
		defer rows.Close()

		type User struct {
			UserID    int64   `json:"user_id"`
			SessionID string  `json:"session_id"`
			Balance   float64 `json:"balance"`
		}

		users := []User{}

		for rows.Next() {
			var u User
			if err := rows.Scan(&u.UserID, &u.SessionID, &u.Balance); err != nil {
				return c.Status(500).SendString("Failed to scan user row")
			}
			users = append(users, u)
		}

		return c.JSON(users)
	})

	//เริ่มจากการสร้าง user แล้วการใส่ learner , teacher , admin ให้ user
	// app.Post("/users/registerUser"       , func(c *fiber.Ctx) error {...}
	// app.Post("/users/:id/registerLearner", func(c *fiber.Ctx) error {...})
	// app.Post("/users/:id/registerTeacher", func(c *fiber.Ctx) error {...})
	// app.Post("/users/:id/registerAdmin"  , func(c *fiber.Ctx) error {...})

	log.Fatal(app.Listen(":3000")) //using PORT 3000 (localhost:3000)
}
