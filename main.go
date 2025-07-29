package main

import (
	"context"
	"log"

	//module name "github.com/a2n2k3p4/tutorium-backend"
	"github.com/a2n2k3p4/tutorium-backend/config/dbserver" //store functions related to connecting to PostgreSQL
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
		// 1) Select all the columns you need
		rows, err := db.Query(context.Background(), `
        	SELECT 
          	user_id,
          	first_name,
          	last_name,
          	gender,
          	phone_number,
          	balance,
          	learner_id,
          	teacher_id,
          	admin_id
        	FROM users
    	`)
		if err != nil {
			return c.Status(500).SendString("Failed to query users")
		}
		defer rows.Close()

		// 2) Define a struct matching the columns (use pointers for nullable FKs)
		type User struct {
			UserID      int64   `json:"user_id"`
			FirstName   string  `json:"first_name"`
			LastName    string  `json:"last_name"`
			Gender      string  `json:"gender"`
			PhoneNumber string  `json:"phone_number"`
			Balance     float64 `json:"balance"`
			LearnerID   *int64  `json:"learner_id,omitempty"`
			TeacherID   *int64  `json:"teacher_id,omitempty"`
			AdminID     *int64  `json:"admin_id,omitempty"`
		}

		users := make([]User, 0)

		// 3) Scan into the matching fields in the same order
		for rows.Next() {
			var u User
			if err := rows.Scan(
				&u.UserID,
				&u.FirstName,
				&u.LastName,
				&u.Gender,
				&u.PhoneNumber,
				&u.Balance,
				&u.LearnerID,
				&u.TeacherID,
				&u.AdminID,
			); err != nil {
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
