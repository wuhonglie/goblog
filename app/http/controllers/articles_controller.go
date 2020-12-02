package controllers

import (
    "fmt"
    article2 "goblog/app/models/article"
    "goblog/pkg/logger"
    "goblog/pkg/route"
    "goblog/pkg/types"
    "gorm.io/gorm"
    "html/template"
    "net/http"
)

type ArticlesController struct {

}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request){
    id := route.GetRouteVariable("id", r)
    article, err := article2.Get(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
            "RouteName2URL": route.Name2URL,
            "Int64ToString": types.Int64ToString,
        }).ParseFiles("resources/views/articles/show.gohtml")
        logger.LogError(err)
        tmpl.Execute(w, article)
    }
}
// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
    articles, err := article2.GetAll()
    //fmt.Println("文章数据：",articles)
    if err != nil {
        logger.LogError(err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, "500 服务器内部错误")
    } else {
        tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
        logger.LogError(err)
        //fmt.Println("文章数据：",articles)
        tmpl.Execute(w, articles)

    }
}
