package models

import (
	"encoding/json"
	"time"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusRejected  BookingStatus = "rejected"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
	StatusCheckedIn BookingStatus = "checked_in"
	StatusProposed  BookingStatus = "proposed"
)

// Booking represents a reservation of a resource
type Booking struct {
	ID              string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ResourceID      string          `json:"resourceId" gorm:"index;type:varchar(36);not null"`
	UserID          string          `json:"userId" gorm:"index;type:varchar(36);not null"`
	Start           time.Time       `json:"start" gorm:"not null"`
	End             time.Time       `json:"end" gorm:"not null"`
	Status          BookingStatus   `json:"status" gorm:"index;type:varchar(20);default:'pending'"`
	CreatedAt       time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	RejectionReason *string         `json:"rejectionReason,omitempty" gorm:"type:text"`
	Details         json.RawMessage `json:"details" gorm:"type:json"` // Stored as JSON
}
