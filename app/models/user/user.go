package user

import (
    "goblog/app/models"
    "goblog/pkg/password"
)

type User struct {
    models.BaseModel
    Name string `gorm:"type:varchar(255);not null;unique" valid:"name"`
    Email string `gorm:"type:varchar(255);unique;" valid:"email"`
    Password string `gorm:"type:varchar(255)" valid:"password"`
    PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (u User) ComparePassword(_password string) bool {
    return password.CheckHash(_password, u.Password)
}
func (u User) Link() string {
    return ""
}
