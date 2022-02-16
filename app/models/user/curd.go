package user

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// Create 创建文章，通过 article.ID 来判断是否创建成功
func (user *User) Create() (err error) {
	result := model.DB.Create(&user)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// GetByEmail 通过 Email 来获取用户
func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get 通过 ID 获取用户
func Get(idstr string) (User, error) {
	var user User
	id := types.StringToInt(idstr)
	if err := model.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	return user, nil
}

// All 获取所有用户数据
func All() ([]User, error) {
	var users []User
	if err := model.DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
