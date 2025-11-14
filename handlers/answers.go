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

type Answer struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	QuestionID int       `json:"question_id" gorm:"column:question_id"`
	UserID     string    `json:"user_id" gorm:"column:user_id"`
	Text       string    `json:"text" gorm:"column:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func MakeAnswersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			handleAnswerByID(db, w, r)

		case http.MethodPost:
			handleAnswerCreate(db, w, r)

		case http.MethodDelete:
			handleAnswerDelete(db, w, r)

		default:
			log.Printf("Method not allowed: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleAnswerCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var body struct {
		QuestionID int    `json:"question_id"`
		UserID     string `json:"user_id"`
		Text       string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "cannot decode body", http.StatusBadRequest)
		return
	}

	log.Printf("CreateAnswer: question_id=%d user_id=%s", body.QuestionID, body.UserID)

	if body.QuestionID == 0 || body.UserID == "" || body.Text == "" {
		log.Println("Validation error: missing required fields")
		http.Error(w, "question_id, user_id and text are required", http.StatusBadRequest)
		return
	}

	var q Question
	if err := db.First(&q, body.QuestionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Question not found for answer create: id=%d", body.QuestionID)
			http.Error(w, "question not found", http.StatusBadRequest)
		} else {
			log.Printf("DB error loading question for answer create: %v", err)
			http.Error(w, "cannot find question(lost db conn?)", http.StatusBadRequest)
		}
		return
	}

	a := Answer{
		QuestionID: body.QuestionID,
		UserID:     body.UserID,
		Text:       body.Text,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&a).Error; err != nil {
		log.Printf("DB error saving answer: %v", err)
		http.Error(w, "failed to save answer", http.StatusInternalServerError)
		return
	}

	log.Printf("Answer created: id=%d for question_id=%d", a.ID, a.QuestionID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

func handleAnswerDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	trimQ := strings.TrimPrefix(r.URL.Path, "/answers/")
	idStr := strings.TrimSuffix(trimQ, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid answer id: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&Answer{}, id).Error; err != nil {
		log.Printf("DB error deleting answer id=%d: %v", id, err)
		http.Error(w, "failed to delete answer", http.StatusInternalServerError)
		return
	}

	log.Printf("Answer deleted: id=%d", id)

	w.WriteHeader(http.StatusOK)
}

func handleAnswerByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	trimQ := strings.TrimPrefix(r.URL.Path, "/answers/")
	idStr := strings.TrimSuffix(trimQ, "/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Printf("Invalid answer id: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var a Answer
	if err := db.First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Answer not found: id=%d", id)
			http.Error(w, "answer not found", http.StatusNotFound)
		} else {
			log.Printf("DB error loading answer id=%d: %v", id, err)
			http.Error(w, "failed to load answer", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Answer loaded: id=%d", id)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(a)
}
