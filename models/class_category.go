package models

import "gorm.io/gorm"

type ClassCategory struct {
	gorm.Model
	ClassCategory string `json:"class_category" gorm:"size:30;unique"`

	Classes  []Class   `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
	Learners []Learner `gorm:"many2many:interested_class_categories;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassCategoryDoc struct {
	ClassCategory string     `json:"class_category" example:"Mathematics"`
	Classes       []ClassDoc `json:"classes,omitempty"`
}

type ClassCategoryIDsDoc struct {
	ClassCategoryIDs []int `json:"class_category_ids" swaggertype:"array,integer" example:"1"`
}

type ClassCategoriesDoc struct {
	Categories []string `json:"categories" swaggertype:"array,string" example:"Mathematics"`
}
