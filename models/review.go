package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	LearnerID uint `gorm:"not null;uniqueIndex:idx_learner_class"`
	ClassID   uint `gorm:"not null;uniqueIndex:idx_learner_class"`
	Rating    int  `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ReviewDoc struct {
	ID        uint   `json:"id" example:"101"`
	LearnerID uint   `json:"learner_id" example:"42"`
	ClassID   uint   `json:"class_id" example:"9"`
	Rating    int    `json:"rating" example:"5"`
	Comment   string `json:"comment" example:"This class was very informative and well-structured."`
}
