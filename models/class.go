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
	Rating           float64         `json:"rating" gorm:"type:numeric(3,2);default:0;check:rating >= 0 AND rating <= 5"`
	Teacher          Teacher         `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
	Categories       []ClassCategory `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
}

//NOTE: You can find class_category relation in CLASS_CLASS_CATEGORY table in Database (Gorm will create it automatically)
//Table will have 2 fields : {class_id, class_category_id}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassDoc struct {
	TeacherID        uint    `json:"teacher_id" example:"7"`
	ClassName        string  `json:"class_name" example:"Advanced Python Programming"`
	ClassDescription string  `json:"class_description" example:"Advanced Python programming course"`
	BannerPicture    string  `json:"banner_picture,omitempty" example:"<base64-encoded-image>"`
	Rating           float64 `json:"rating" example:"4.7"`
}
