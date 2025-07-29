package main

import (
	"log"

	//module name "github.com/Parkorn/KUTutorium"
	"github.com/Parkorn/KUTutorium/config/dbserver" //store functions related to connecting to PostgreSQL
	"github.com/gofiber/fiber/v2"
)

// Before running the server, change config/dbserver/config.go to correct connection info
// Define a struct matching the columns (use pointers for nullable FKs)
type User struct {
	UserID      int64   `json:"user_id" gorm:"column:user_id"`
	FirstName   string  `json:"first_name" gorm:"column:first_name"`
	LastName    string  `json:"last_name" gorm:"column:last_name"`
	Gender      string  `json:"gender" gorm:"column:gender"`
	PhoneNumber string  `json:"phone_number" gorm:"column:phone_number"`
	Balance     float64 `json:"balance" gorm:"column:balance"`
	LearnerID   *int64  `json:"learner_id,omitempty" gorm:"column:learner_id"`
	TeacherID   *int64  `json:"teacher_id,omitempty" gorm:"column:teacher_id"`
	AdminID     *int64  `json:"admin_id,omitempty" gorm:"column:admin_id"`
}

func main() {
	cfg := dbserver.NewConfig()

	db, err := dbserver.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	app := fiber.New()

	// Define the /users route and handler inline
	app.Get("/users", func(c *fiber.Ctx) error {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			return c.Status(500).SendString("Failed to query users")
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
