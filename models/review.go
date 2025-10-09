package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	LearnerID uint   `json:"learner_id" gorm:"uniqueIndex:idx_learner_class"`
	ClassID   uint   `json:"class_id" gorm:"not null;uniqueIndex:idx_learner_class"`
	Rating    int    `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string `json:"comment"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:SET NULL"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ReviewDoc struct {
	LearnerID uint   `json:"learner_id" example:"42"`
	ClassID   uint   `json:"class_id" example:"9"`
	Rating    int    `json:"rating" example:"5"`
	Comment   string `json:"comment" example:"This class was very informative and well-structured."`
}
