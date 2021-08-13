package user

/**
GORM 模型钩子 是在创建、查询、更新、删除等操作之前、之后调用的函数。
为模型定义指定的方法，它会在创建、更新、查询、删除时自动被调用。
如果任何回调返回错误，GORM 将停止后续的操作并回滚事务。
*/
import (
	"goblog/pkg/password"

	"gorm.io/gorm"
)

//// BeforeCreate GORM 的模型钩子，创建模型前调用
//func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
//	u.Password = password.Hash(u.Password)
//	return
//}
//
//// BeforeUpdate GORM 的模型钩子，更新模型前调用
//func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
//	if !password.IsHashed(u.Password) {
//		u.Password = password.Hash(u.Password)
//	}
//	return
//}

// BeforeSave GORM 的模型钩子，在保存和更新模型前调用
//针对创建和更新的事件监控，我们可以使用覆盖面更广的 BeforeSave 钩子
func (u *User) BeforeSave(tx *gorm.DB) (err error) {

	if !password.IsHashed(u.Password) {
		u.Password = password.Hash(u.Password)
	}
	return
}
