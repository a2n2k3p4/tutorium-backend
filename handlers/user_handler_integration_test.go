package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_UserCRUD(t *testing.T) {
	updatedPhone := "+66111111111"
	updatedBanCount := 1

	runCRUDTest(t, crudTestCase[models.User]{
		ResourceName: "users",
		BasePath:     "/users/",
		Create: func(t *testing.T) models.User {
			user, _ := createTestUser(t)
			return user
		},
		GetID: func(u models.User) uint { return u.ID },
		UpdatePayload: map[string]any{
			"phone_number": updatedPhone,
			"ban_count":    updatedBanCount,
		},
		AssertUpdated: func(t *testing.T, updated models.User) {
			if updated.PhoneNumber != updatedPhone || updated.BanCount != updatedBanCount {
				t.Fatalf("update failed, expected phone=%s ban=%d got phone=%s ban=%d", updatedPhone, updatedBanCount, updated.PhoneNumber, updated.BanCount)
			}
		},
	})
}
