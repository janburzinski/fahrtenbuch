package util

import (
	"testing"
)

func TestGenerateAndVerify(t *testing.T) {
	randomPassword, err := GenerateRandomBytes(8)
	if err != nil {
		t.Fatalf("error while generating random password: %s", err)
		return
	}

	p := &Argon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	hashedPassword, err := GenerateFromPassword(string(randomPassword), p)
	if err != nil {
		t.Fatalf("error while generating password hash: %s", err)
		return
	}

	match, err := ComparePasswordAndHash(string(randomPassword), hashedPassword)
	if err != nil {
		t.Fatalf("error while verifying password: %s", err)
		return
	}

	if !match {
		t.Fatal("password and password hash dont match")
		return
	}
}
