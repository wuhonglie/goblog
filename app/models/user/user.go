package user

import "goblog/app/models"

type User struct {
    models.BaseModel
    Name string `gorm:"type:varchar(255);not null;unique" valid:"name"`
    Email string `gorm:"type:varchar(255);unique;" valid:"email"`
    Password string `gorm:"type:varchar(255)" valid:"password"`
    PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (u User) ComparePassword(password string) bool {
    return u.Password == password
}
