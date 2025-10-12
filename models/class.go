package models

import (
	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	TeacherID        uint            `json:"teacher_id" gorm:"not null"`
	ClassName        string          `json:"class_name" gorm:"size:255;not null"`
	ClassDescription string          `json:"class_description" gorm:"size:1000"`
	BannerPictureURL string          `json:"banner_picture,omitempty"`
	Teacher          Teacher         `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
	Categories       []ClassCategory `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassDoc struct {
	TeacherID        uint   `json:"teacher_id" example:"7"`
	ClassName        string `json:"class_name" example:"Advanced Python Programming"`
	ClassDescription string `json:"class_description" example:"Advanced Python programming course"`
	BannerPicture    string `json:"banner_picture,omitempty" example:"<base64-encoded-image>"`
}

type ClassAverageRating struct {
	ClassID       uint    `json:"class_id,omitempty" example:"123"`
	AverageRating float64 `json:"average_rating" example:"4.5"`
}
