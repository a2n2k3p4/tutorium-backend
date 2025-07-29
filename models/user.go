package models

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
