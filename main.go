package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Config struct {
	QuizUrl string `yaml:"quizUrl"`
}

type Quiz struct {
	QuizId    string     `yaml:"quizId"`
	Questions []Question `yaml:"questions"`
}

type Question struct {
	ID             int      `yaml:"id"`
	Question       string   `yaml:"question"`
	Type           string   `yaml:"type"` // "single" or "multiple"
	Options        []string `yaml:"options"`
	CorrectAnswers []int    `yaml:"answers"`
}

type EvaluationResult struct {
	QuestionID int  `json:"questionId"`
	Correct    bool `json:"correct"`
}

type Submission struct {
	QuizId           string              `json:"quizId"`
	SubmittedAnswers map[string][]string `json:"answers"`
}

var config Config

func init() {
	// Load config.yaml
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/check", evaluateHandler)
	http.HandleFunc("/health", healthCheckHandler)
	log.Println("Starting external service on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All is well here")
}

func evaluateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission Submission

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &submission)
	if err != nil {
		http.Error(w, "Failed to unmarshal payload", http.StatusBadRequest)
		return
	}

	quizId := submission.QuizId

	quiz, err := fetchQuizFromExternalService(quizId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load quiz: %v", err), http.StatusInternalServerError)
		return
	}

	evaluationResults := evaluateAnswers(quiz, submission.SubmittedAnswers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evaluationResults)
}

func fetchQuizFromExternalService(quizID string) (*Quiz, error) {
	url := fmt.Sprintf("%s%s", config.QuizUrl, quizID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quiz from external service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var quiz Quiz
	err = yaml.Unmarshal(body, &quiz)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	return &quiz, nil
}

func evaluateAnswers(quiz *Quiz, submittedAnswers map[string][]string) []EvaluationResult {
	var results []EvaluationResult

	for _, question := range quiz.Questions {
		submitted := submittedAnswers[fmt.Sprintf("question-%d", question.ID)]

		// Convert submitted answers to integers
		var submittedInt []int
		for _, s := range submitted {
			var temp int
			fmt.Sscan(s, &temp)
			submittedInt = append(submittedInt, temp)
		}

		// Check if the submitted answers match the correct answers
		correct := reflect.DeepEqual(submittedInt, question.CorrectAnswers)

		results = append(results, EvaluationResult{
			QuestionID: question.ID,
			Correct:    correct,
		})
	}

	return results
}
