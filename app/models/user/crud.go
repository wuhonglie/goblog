package user

import (
    "goblog/pkg/logger"
    "goblog/pkg/model"
    "goblog/pkg/types"
)

func (user *User) Create() (err error) {
    if err = model.DB.Create(&user).Error; err != nil {
        logger.LogError(err)
        return err
    }
    return nil
}

func Get(idstr string) (User,error) {
    var user User
    id := types.StringToInt(idstr)
    if err := model.DB.First(&user,id).Error; err != nil {
        return user, err
    }
    return user, nil
}

func GetByEmail(email string) (User,error) {
    var user User
    if result := model.DB.Where("email = ?",email).First(&user); result.Error != nil {
        return user, result.Error
    }
    return user, nil
}
