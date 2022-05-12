package main

import (
	"errors"
	"testing"
)

func TestUserIsEligible(t *testing.T) {
	var tests = []struct {
		email       string
		password    string
		age         int
		expectedErr error
	}{
		{
			email:       "test@example.com",
			password:    "12345",
			age:         18,
			expectedErr: nil,
		},
		{
			email:       "",
			password:    "12345",
			age:         18,
			expectedErr: errors.New("email can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "",
			age:         18,
			expectedErr: errors.New("password can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "12345",
			age:         16,
			expectedErr: errors.New("age must be at least 18"),
		},
	}

	for _, test := range tests {
		err := userIsEligible(test.email, test.password, test.age)

		errString := ""
		expectedErrString := ""

		if err != nil {
			errString = err.Error()
		}
		if test.expectedErr != nil {
			expectedErrString = test.expectedErr.Error()
		}

		if errString != expectedErrString {
			t.Errorf("got %s, want %s", errString, expectedErrString)
		}
	}
}
