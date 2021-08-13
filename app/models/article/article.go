package article

import (
	"goblog/app/models"
	"goblog/pkg/route"
	"strconv"
)

type Article struct {
	models.BaseModel
	Title string `gorm:"type:varchar(255);not null" valid:"title"`
	Body  string `gorm:"type:varchar(255);not null" valid:"body"`
}

/**
type Object struct {
	...
}
// Object 的方法
func (obj *Object) method() {
	...
}

// 只是一个函数
func function() {
	...
}

调用的对比
// 调用方法：
o := new(Object)
o.method()

// 调用函数
function()
*/

// Link 方法用来生成文章链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(int64(article.ID), 10))
}
