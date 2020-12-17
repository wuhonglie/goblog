package article

import (
    "goblog/app/models"
)

type Article struct {
    models.BaseModel
    //ID int64
    Title string
    Body string
}
//func (a Article) Link() string {
//    return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
//}
//.
//.
//.
// Link 方法用来生成文章链接
//func (a Article) Link() string {
//   return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
//}
