package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestGenerateAndVerifyJWT(t *testing.T) {
	// run the test with the accessToken fail enabled and disabled
	var tests = []struct {
		access bool
	}{
		{true},
		{false},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("accessToken? %t", tt.access)
		t.Run(testName, func(t *testing.T) {
			userId := strconv.Itoa(randomUserId())
			token, err := CreateToken(userId, tt.access)
			if err != nil {
				t.Logf("error while creating token with %t, err: %s", tt.access, err.Error())
				t.Fail()
			}

			err = VerifyToken(token, tt.access)
			if err != nil {
				t.Logf("error while verifying access token: %s", err.Error())
				t.Fail()
			}
		})
	}
}

func TestGetUserIdFromJWT(t *testing.T) {
	id := strconv.Itoa(randomUserId())
	token, err := CreateToken(id, true)
	if err != nil {
		t.Logf("error while creating token: %s", err.Error())
		t.Fail()
	}

	userId, err := GetUserIdFromJWT(token)
	if err != nil {
		t.Logf("error while getting user id from jwt: %s", err.Error())
		t.Fail()
	}

	if userId != id {
		t.Logf("expected: %s, got: %s", id, userId)
		t.Fail()
	}
}

func randomUserId() int {
	return rand.Int()
}
