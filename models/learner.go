package models

import "gorm.io/gorm"

type Learner struct {
	gorm.Model
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type LearnerDoc struct {
	ID uint `json:"id" example:"42"`
}
