package dto

import (
	"fmt"
	"subscription/models"
	"time"
)

// --- Request DTO ---
type SubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // формат "MM-YYYY"
	EndDate     string `json:"end_date,omitempty"`            // формат "MM-YYYY"
}

// Конвертация в модель
func (r SubscriptionRequest) ToModel() (models.Subscription, error) {
	// парсим start_date (MM-YYYY → time.Time)
	start, err := time.Parse("01-2006", r.StartDate)
	if err != nil {
		return models.Subscription{}, fmt.Errorf("invalid start_date: %w", err)
	}

	var end *time.Time
	if r.EndDate != "" {
		e, err := time.Parse("01-2006", r.EndDate)
		if err != nil {
			return models.Subscription{}, fmt.Errorf("invalid end_date: %w", err)
		}
		end = &e
	}

	return models.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      r.UserID,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

// --- Response DTO ---
type SubscriptionResponse struct {
	ID          uint   `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`         // формат "MM-YYYY"
	EndDate     string `json:"end_date,omitempty"` //формат "MM-YYYY"
}

// Конвертация из модели
func FromModel(m models.Subscription) SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:          m.ID,
		ServiceName: m.ServiceName,
		Price:       m.Price,
		UserID:      m.UserID,
		StartDate:   m.StartDate.Format("01-2006"), // форматируем обратно
	}
	if m.EndDate != nil {
		resp.EndDate = m.EndDate.Format("01-2006")
	}
	return resp
}
