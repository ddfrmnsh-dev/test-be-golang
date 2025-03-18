package model

import "time"

type User struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null;size:255"`
	Email     string    `json:"email" gorm:"not null;size:255"`
	Password  string    `json:"password" gorm:"not null;size:255"`
	Role      string    `json:"role" gorm:"not null;size:50"`
	IsActive  *bool     `json:"isActive" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type GetCustomerDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}

type InputLogin struct {
	Identifier string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type InputRegister struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	IsActive *bool  `json:"isActive" binding:"required"`
}

type UserResponse struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	IsActive  *bool     `json:"isActive"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FormatUserResponse(user User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Username:  user.Username,
		IsActive:  user.IsActive,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
