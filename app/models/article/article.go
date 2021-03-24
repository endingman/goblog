package article

import (
	"goblog/pkg/route"
	"strconv"
)

type Article struct {
	ID    int
	Title string
	Body  string
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
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(int64(a.ID), 10))
}
