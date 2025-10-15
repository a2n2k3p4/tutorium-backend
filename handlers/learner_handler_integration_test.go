package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Learner_CRUD(t *testing.T) {
	runCRUDTest(t, crudTestCase[models.Learner]{
		ResourceName: "learners",
		BasePath:     "/learners/",
		Create: func(t *testing.T) models.Learner {
			_, learner := createTestUser(t)
			return learner
		},
		GetID: func(l models.Learner) uint { return l.ID },
	})
}

func TestIntegration_Learner_Interests(t *testing.T) {
	_, learner := createTestUser(t)
	catA := createTestClassCategory(t)
	catB := createTestClassCategory(t)

	addPayload := map[string]any{
		"class_category_ids": []uint{catA.ID, catB.ID},
	}
	var updated models.Learner
	jsonRequestExpect(
		t,
		http.MethodPost,
		fmt.Sprintf("/learners/%d/interests", learner.ID),
		addPayload,
		http.StatusOK,
		&updated,
	)
	if len(updated.Interested) != 2 {
		t.Fatalf("expected 2 interests, got %d", len(updated.Interested))
	}

	var interestResp struct {
		Categories []string `json:"categories"`
	}
	jsonRequestExpect(
		t,
		http.MethodGet,
		fmt.Sprintf("/learners/%d/interests", learner.ID),
		nil,
		http.StatusOK,
		&interestResp,
	)
	if len(interestResp.Categories) != 2 {
		t.Fatalf("expected 2 category names, got %d", len(interestResp.Categories))
	}
	sort.Strings(interestResp.Categories)
	wantCats := []string{catA.ClassCategory, catB.ClassCategory}
	sort.Strings(wantCats)
	for i, name := range wantCats {
		if interestResp.Categories[i] != name {
			t.Fatalf("expected categories %v, got %v", wantCats, interestResp.Categories)
		}
	}

	removePayload := map[string]any{
		"class_category_ids": []uint{catA.ID},
	}
	jsonRequestExpect(
		t,
		http.MethodDelete,
		fmt.Sprintf("/learners/%d/interests", learner.ID),
		removePayload,
		http.StatusOK,
		&updated,
	)
	if len(updated.Interested) != 1 || updated.Interested[0].ID != catB.ID {
		t.Fatalf("expected only category %d to remain, got %+v", catB.ID, updated.Interested)
	}
}

func TestIntegration_Learner_Recommended(t *testing.T) {
	user, learner := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)

	catMatch := createTestClassCategory(t)
	catOther := createTestClassCategory(t)
	catNoClass := createTestClassCategory(t)

	classMatch := createTestClass(t, teacher.ID)
	classOther := createTestClass(t, teacher.ID)

	createTestClassSession(t, classMatch.ID)
	createTestClassSession(t, classOther.ID)

	var updated models.Learner

	var noInterestResp struct {
		RecommendedFound   bool             `json:"recommended_found"`
		RecommendedClasses []models.Class   `json:"recommended_classes"`
		RemainingClasses   []models.Class   `json:"remaining_classes"`
	}
	jsonRequestExpect(
		t,
		http.MethodGet,
		fmt.Sprintf("/learners/%d/recommended", learner.ID),
		nil,
		http.StatusOK,
		&noInterestResp,
	)
	if noInterestResp.RecommendedFound {
		t.Fatalf("expected no interests to yield RecommendedFound=false")
	}
	if len(noInterestResp.RecommendedClasses) != 0 {
		t.Fatalf("expected no recommended classes, got %d", len(noInterestResp.RecommendedClasses))
	}
	if len(noInterestResp.RemainingClasses) == 0 {
		t.Fatalf("expected remaining classes to be non-empty when learner has no interests")
	}

	jsonRequestExpect(
		t,
		http.MethodPost,
		fmt.Sprintf("/classes/%d/categories", classMatch.ID),
		map[string]any{"class_category_ids": []uint{catMatch.ID}},
		http.StatusOK,
		nil,
	)
	jsonRequestExpect(
		t,
		http.MethodPost,
		fmt.Sprintf("/classes/%d/categories", classOther.ID),
		map[string]any{"class_category_ids": []uint{catOther.ID}},
		http.StatusOK,
		nil,
	)

	jsonRequestExpect(
		t,
		http.MethodPost,
		fmt.Sprintf("/learners/%d/interests", learner.ID),
		map[string]any{"class_category_ids": []uint{catMatch.ID}},
		http.StatusOK,
		nil,
	)

	var resp struct {
		RecommendedFound   bool           `json:"recommended_found"`
		RecommendedClasses []models.Class `json:"recommended_classes"`
		RemainingClasses   []models.Class `json:"remaining_classes"`
	}

	jsonRequestExpect(
		t,
		http.MethodGet,
		fmt.Sprintf("/learners/%d/recommended", learner.ID),
		nil,
		http.StatusOK,
		&resp,
	)

	if !resp.RecommendedFound {
		t.Fatalf("expected recommended classes to be found")
	}
	if len(resp.RecommendedClasses) != 1 || resp.RecommendedClasses[0].ID != classMatch.ID {
		t.Fatalf("expected recommended class %d, got %+v", classMatch.ID, resp.RecommendedClasses)
	}
	if len(resp.RemainingClasses) == 0 {
		t.Fatalf("expected remaining classes to be non-empty")
	}
	foundOther := false
	for _, cls := range resp.RemainingClasses {
		if cls.ID == classOther.ID {
			foundOther = true
			break
		}
	}
	if !foundOther {
		t.Fatalf("expected class %d to appear in remaining classes", classOther.ID)
	}

	jsonRequestExpect(
		t,
		http.MethodDelete,
		fmt.Sprintf("/learners/%d/interests", learner.ID),
		map[string]any{"class_category_ids": []uint{catMatch.ID}},
		http.StatusOK,
		&updated,
	)

	jsonRequestExpect(
		t,
		http.MethodPost,
		fmt.Sprintf("/learners/%d/interests", learner.ID),
		map[string]any{"class_category_ids": []uint{catNoClass.ID}},
		http.StatusOK,
		nil,
	)

	var noMatchResp struct {
		RecommendedFound   bool             `json:"recommended_found"`
		RecommendedClasses []models.Class   `json:"recommended_classes"`
		RemainingClasses   []models.Class   `json:"remaining_classes"`
	}
	jsonRequestExpect(
		t,
		http.MethodGet,
		fmt.Sprintf("/learners/%d/recommended", learner.ID),
		nil,
		http.StatusOK,
		&noMatchResp,
	)
	if noMatchResp.RecommendedFound {
		t.Fatalf("expected RecommendedFound=false when no classes match interests")
	}
	if len(noMatchResp.RecommendedClasses) != 0 {
		t.Fatalf("expected no recommended classes when interests have no matches, got %d", len(noMatchResp.RecommendedClasses))
	}
	if len(noMatchResp.RemainingClasses) == 0 {
		t.Fatalf("expected remaining classes to be populated when no matches are found")
	}
}
