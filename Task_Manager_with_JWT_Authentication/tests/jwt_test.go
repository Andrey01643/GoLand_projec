package tests

import (
	jwt2 "CRUD_REST/jwt"
	"testing"
)

func TestJWTToken(t *testing.T) {
	userID := 1
	token, err := jwt2.GenerateJWTToken(userID)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}

	validUserID, err := jwt2.ValidateJWTToken(token)
	if err != nil {
		t.Fatalf("failed to validate JWT token: %v", err)
	}
	if validUserID != userID {
		t.Fatalf("expected userID %d, but got %d", userID, validUserID)
	}

	invalidToken := "invalid_token"
	_, err = jwt2.ValidateJWTToken(invalidToken)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
	expectedErrorMsg := "invalid token"
	if err.Error() != expectedErrorMsg {
		t.Fatalf("expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
	}

	invalidUserID := 0
	_, err = jwt2.GenerateJWTToken(invalidUserID)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
	expectedErrorMsg = "invalid userID"
	if err.Error() != expectedErrorMsg {
		t.Fatalf("expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
	}
}
