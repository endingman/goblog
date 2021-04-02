package user

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
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
