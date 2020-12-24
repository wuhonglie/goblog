package article

import (
    "goblog/app/models"
    "goblog/app/models/user"
    "goblog/pkg/route"
)

type Article struct {
    models.BaseModel
    //ID int64
    Title string `valid:"title"`
    Body string `valid:"body"`
    UserID uint64 `gorm:"not null;index"`
    User user.User
}
//func (a Article) Link() string {
//    return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
//}
//.
//.
//.
// Link 方法用来生成文章链接
func (a Article) Link() string {
 return route.Name2URL("articles.show", "id", a.GetStringID())
}

func (a Article) CreatedAtDate() string {
    return a.CreatedAt.Format("2006-01-02")
}
