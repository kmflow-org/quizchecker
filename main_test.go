package main

import (
	"reflect"
	"testing"
)

// Unit test for evaluateAnswers
func TestEvaluateAnswers(t *testing.T) {
	// Test cases
	tests := []struct {
		name             string
		quiz             Quiz
		submittedAnswers map[string][]string
		expectedResults  []EvaluationResult
	}{
		{
			name: "All answers correct",
			quiz: Quiz{
				Questions: []Question{
					{ID: 1, CorrectAnswers: []int{1, 2}},
					{ID: 2, CorrectAnswers: []int{3}},
				},
			},
			submittedAnswers: map[string][]string{
				"question-1": {"1", "2"},
				"question-2": {"3"},
			},
			expectedResults: []EvaluationResult{
				{QuestionID: 1, Correct: true},
				{QuestionID: 2, Correct: true},
			},
		},
		{
			name: "Some answers incorrect",
			quiz: Quiz{
				Questions: []Question{
					{ID: 1, CorrectAnswers: []int{1, 2}},
					{ID: 2, CorrectAnswers: []int{3}},
				},
			},
			submittedAnswers: map[string][]string{
				"question-1": {"1", "3"},
				"question-2": {"4"},
			},
			expectedResults: []EvaluationResult{
				{QuestionID: 1, Correct: false},
				{QuestionID: 2, Correct: false},
			},
		},
		{
			name: "Empty submitted answers",
			quiz: Quiz{
				Questions: []Question{
					{ID: 1, CorrectAnswers: []int{1, 2}},
					{ID: 2, CorrectAnswers: []int{3}},
				},
			},
			submittedAnswers: map[string][]string{
				"question-1": {},
				"question-2": {},
			},
			expectedResults: []EvaluationResult{
				{QuestionID: 1, Correct: false},
				{QuestionID: 2, Correct: false},
			},
		},
		{
			name: "Partially correct answers",
			quiz: Quiz{
				Questions: []Question{
					{ID: 1, CorrectAnswers: []int{1, 2}},
					{ID: 2, CorrectAnswers: []int{3}},
				},
			},
			submittedAnswers: map[string][]string{
				"question-1": {"1", "2"},
				"question-2": {"4"},
			},
			expectedResults: []EvaluationResult{
				{QuestionID: 1, Correct: true},
				{QuestionID: 2, Correct: false},
			},
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := evaluateAnswers(&tt.quiz, tt.submittedAnswers)
			if !reflect.DeepEqual(results, tt.expectedResults) {
				t.Errorf("evaluateAnswers() = %v, want %v", results, tt.expectedResults)
			}
		})
	}
}
