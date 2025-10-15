package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Class_CRUD(t *testing.T) {
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	updatedDescription := "Updated integration class description"

	runCRUDTest(t, crudTestCase[models.Class]{
		ResourceName: "classes",
		BasePath:     "/classes/",
		Create: func(t *testing.T) models.Class {
			return createTestClass(t, teacher.ID)
		},
		GetID: func(c models.Class) uint { return c.ID },
		UpdatePayload: map[string]any{
			"class_description": updatedDescription,
		},
		AssertUpdated: func(t *testing.T, updated models.Class) {
			if updated.ClassDescription != updatedDescription {
				t.Fatalf("expected updated description %q , got desc=%q ", updatedDescription, updated.ClassDescription)
			}
		},
	})
}

func TestIntegration_Class_CategoriesRoutes(t *testing.T) {
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	class := createTestClass(t, teacher.ID)

	cat1 := createTestClassCategory(t)
	cat2 := createTestClassCategory(t)

	var updated models.Class
	jsonRequestExpect(
		t,
		http.MethodPost,
		fmt.Sprintf("/classes/%d/categories", class.ID),
		map[string]any{"class_category_ids": []uint{cat1.ID, cat2.ID}},
		http.StatusOK,
		&updated,
	)
	assertHasCategory := func(cats []models.ClassCategory, id uint) bool {
		for _, cat := range cats {
			if cat.ID == id {
				return true
			}
		}
		return false
	}
	if !assertHasCategory(updated.Categories, cat1.ID) || !assertHasCategory(updated.Categories, cat2.ID) {
		t.Fatalf("expected categories %d and %d to be attached, got %+v", cat1.ID, cat2.ID, updated.Categories)
	}

	var listResp struct {
		Categories []string `json:"categories"`
	}
	jsonRequestExpect(
		t,
		http.MethodGet,
		fmt.Sprintf("/classes/%d/categories", class.ID),
		nil,
		http.StatusOK,
		&listResp,
	)
	contains := func(list []string, target string) bool {
		for _, v := range list {
			if v == target {
				return true
			}
		}
		return false
	}
	if !contains(listResp.Categories, cat1.ClassCategory) || !contains(listResp.Categories, cat2.ClassCategory) {
		t.Fatalf("expected GET categories to include %q and %q, got %v", cat1.ClassCategory, cat2.ClassCategory, listResp.Categories)
	}

	jsonRequestExpect(
		t,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d/categories", class.ID),
		map[string]any{"class_category_ids": []uint{cat1.ID}},
		http.StatusOK,
		&updated,
	)
	if !assertHasCategory(updated.Categories, cat2.ID) {
		t.Fatalf("expected category %d to remain after delete, got %+v", cat2.ID, updated.Categories)
	}
}
