package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Admin_CRUD(t *testing.T) {
	user, _ := createTestUser(t)

	runCRUDTest(t, crudTestCase[models.Admin]{
		ResourceName: "admins",
		BasePath:     "/admins/",
		Create: func(t *testing.T) models.Admin {
			admin := createTestAdmin(t, user.ID)
			return admin
		},
		GetID: func(a models.Admin) uint { return a.ID },
	})
}
