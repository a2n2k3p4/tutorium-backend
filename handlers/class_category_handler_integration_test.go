package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_ClassCategory_CRUD(t *testing.T) {
	updatedCategory := "Updated Integration Category"

	runCRUDTest(t, crudTestCase[models.ClassCategory]{
		ResourceName: "class categories",
		BasePath:     "/class_categories/",
		Create: func(t *testing.T) models.ClassCategory {
			return createTestClassCategory(t)
		},
		GetID:         func(c models.ClassCategory) uint { return c.ID },
		UpdatePayload: map[string]any{"class_category": updatedCategory},
		AssertUpdated: func(t *testing.T, updated models.ClassCategory) {
			if updated.ClassCategory != updatedCategory {
				t.Fatalf("expected updated category name %q, got %q", updatedCategory, updated.ClassCategory)
			}
		},
	})
}
