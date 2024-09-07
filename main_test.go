package main

import (
	"testing"
)

// Example function to test (replace with actual functions)
func ValidateAnswer(submission string, correctAnswer string) bool {
	return submission == correctAnswer
}

// Unit test for ValidateAnswer function
func TestValidateAnswer(t *testing.T) {
	tests := []struct {
		submission    string
		correctAnswer string
		expected      bool
	}{
		{"A", "A", true},
		{"B", "A", false},
		{"", "A", false},
	}

	for _, test := range tests {
		result := ValidateAnswer(test.submission, test.correctAnswer)
		if result != test.expected {
			t.Errorf("ValidateAnswer(%s, %s) = %v; want %v", test.submission, test.correctAnswer, result, test.expected)
		}
	}
}
