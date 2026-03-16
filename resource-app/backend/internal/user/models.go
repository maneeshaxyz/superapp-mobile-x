package user

import "time"

type User struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Email      string    `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Role       Role      `json:"role" gorm:"type:varchar(20);default:'USER'"`
	Avatar     string    `json:"avatar" gorm:"type:varchar(255)"`
	Department string    `json:"department" gorm:"type:varchar(100)"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
