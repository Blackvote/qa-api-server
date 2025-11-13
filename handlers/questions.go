package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type QuestionWithAnswers struct {
	Question Question `json:"question"`
	Answers  []Answer `json:"answers"`
}

type Question struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func MakeQuestionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			if r.URL.Path == "/questions/" {
				handleQuestionsList(db, w, r)
				return
			}
			handleQuestionByID(db, w, r)

		case "POST":
			if strings.HasSuffix(r.URL.Path, "/answers/") {
				handleCreateAnswerForQuestion(db, w, r)
				return
			}
			handleQuestionCreate(db, w, r)

		case "DELETE":
			handleQuestionDelete(db, w, r)

		default:
			log.Printf("Method not allowed: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleCreateAnswerForQuestion(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	trimQ := strings.TrimPrefix(r.URL.Path, "/questions/")
	idStr := strings.TrimSuffix(trimQ, "/answers/")
	idStr = strings.TrimSuffix(idStr, "/")

	questionID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid question id: %s", idStr)
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	var q Question
	if err := db.First(&q, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Question not found: id=%d", questionID)
			http.Error(w, "question not found", http.StatusNotFound)
		} else {
			log.Printf("DB error loading question: %v", err)
			http.Error(w, "failed to load question", http.StatusInternalServerError)
		}
		return
	}

	var body struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "cannot decode body", http.StatusBadRequest)
		return
	}

	if body.UserID == "" || body.Text == "" {
		log.Printf("Validation error: user_id or text missing")
		http.Error(w, "user_id and text are required", http.StatusBadRequest)
		return
	}

	a := Answer{
		QuestionID: questionID,
		UserID:     body.UserID,
		Text:       body.Text,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&a).Error; err != nil {
		log.Printf("DB error creating answer: %v", err)
		http.Error(w, "failed to add answer", http.StatusInternalServerError)
		return
	}

	log.Printf("Answer created: id=%d for question=%d", a.ID, questionID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

func handleQuestionCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var body struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "cannot decode body", http.StatusBadRequest)
		return
	}

	q := Question{
		Text:      body.Text,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&q).Error; err != nil {
		log.Printf("DB error saving question: %v", err)
		http.Error(w, "failed to save question", http.StatusInternalServerError)
		return
	}

	log.Printf("Question created: id=%d", q.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(q)
}

func handleQuestionsList(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var qs []Question
	if err := db.Find(&qs).Error; err != nil {
		log.Printf("DB error loading questions: %v", err)
		http.Error(w, "failed to load questions", http.StatusInternalServerError)
		return
	}

	log.Printf("Questions loaded: count=%d", len(qs))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(qs)
}

func handleQuestionDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	trimQ := strings.TrimPrefix(r.URL.Path, "/questions/")
	idStr := strings.TrimSuffix(trimQ, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Printf("Invalid id: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var q Question
	if err := db.First(&q, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Question not found: id=%d", id)
			http.Error(w, "question not found", http.StatusNotFound)
		} else {
			log.Printf("DB error loading question: %v", err)
			http.Error(w, "failed to load question", http.StatusInternalServerError)
		}
		return
	}

	if err := db.Delete(&Question{}, id).Error; err != nil {
		log.Printf("DB error deleting question: %v", err)
		http.Error(w, "failed to delete question", http.StatusInternalServerError)
		return
	}

	log.Printf("Question deleted: id=%d", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(q)
}

func handleQuestionByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	trimQ := strings.TrimPrefix(r.URL.Path, "/questions/")
	idStr := strings.TrimSuffix(trimQ, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Printf("Invalid id: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var q Question
	if err := db.First(&q, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Question not found: id=%d", id)
			http.Error(w, "question not found", http.StatusNotFound)
		} else {
			log.Printf("DB error loading question: %v", err)
			http.Error(w, "failed to load question", http.StatusInternalServerError)
		}
		return
	}

	var ans []Answer
	if err := db.Where("question_id = ?", id).Find(&ans).Error; err != nil {
		log.Printf("DB error loading answers: %v", err)
		http.Error(w, "failed to load answers", http.StatusInternalServerError)
		return
	}

	log.Printf("Question loaded: id=%d, answers=%d", id, len(ans))

	resp := QuestionWithAnswers{
		Question: q,
		Answers:  ans,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
