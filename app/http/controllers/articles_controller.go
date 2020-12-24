package controllers

import (
    "fmt"
    article2 "goblog/app/models/article"
    "goblog/app/requests"
    "goblog/pkg/logger"
    "goblog/pkg/route"
    "goblog/pkg/view"
    "gorm.io/gorm"
    "net/http"
)

type ArticlesController struct {

}
//type ArticlesFormData struct {
//    Title, Body string
//    Article article2.Article
//    Errors map[string]string
//}
// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request){
    view.Render(w, view.D{}, "articles.create","articles._form_field")
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
        view.Render(w,  view.D{
            "Article": article,
        },"articles.show")
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
        view.Render(w,  view.D{
            "Articles": articles,
        },"articles.index")

    }
}

//func validateArticleFormData(title string, body string) map[string]string {
//    errors := make(map[string]string)
//    if title == "" {
//        errors["title"] = "标题不能为空"
//    }else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
//        errors["title"] = "标题长度需介于 3-40"
//    }
//    if body == "" {
//        errors["body"] = "内容不能为空"
//    } else if utf8.RuneCountInString(body) < 10 {
//        errors["body"] = "内容长度需大于或等于 10 个字节"
//    }
//    return errors
//}
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request){
    //title := r.PostFormValue("title")
    //body := r.PostFormValue("body")
    // 1.初始化数据
    _article := article2.Article{
        Title: r.PostFormValue("title"),
        Body: r.PostFormValue("body"),
    }
    errors := requests.ValidateArticleForm(_article)
    if len(errors) == 0 {
        _article.Create()
        //lastInsertID, err := saveArticleToDB(title, body)
        if _article.ID > 0 {
            indexURL := route.Name2URL("articles.show","id",_article.GetStringID())
            http.Redirect(w,r,indexURL,http.StatusFound)
        } else {
            //logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "创建文章失败，请联系管理员")
        }
    } else {
        view.Render(w, view.D{
            "Article": _article,
            "Errors": errors,
        }, "articles.create", "articles._form_field")
    }
}

func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article2.Get(id)
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
        view.Render(w, view.D{
            "Article": _article,
            "Errors": view.D{},
        },"articles.edit", "articles._form_field")
    }
}

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article2.Get(id)
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
        _article.Title = r.PostFormValue("title")
        _article.Body = r.PostFormValue("body")
        errors := requests.ValidateArticleForm(_article)
        if len(errors) == 0 {
            rowsAffected, err := _article.Update()
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
                return
            }
            if rowsAffected > 0 {
                showURL := route.Name2URL("articles.show","id",id)
                http.Redirect(w, r, showURL, http.StatusFound)
            } else {
                fmt.Fprint(w, "您没有做任何更改！")
            }
        } else {
            view.Render(w, view.D{
                "Article": _article,
                "Errors": errors,
            },"articles.edit","articles._form_field")
        }
    }
}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article2.Get(id)
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
        rowsAffected, err := _article.Delete()
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        } else {
            if rowsAffected > 0 {
                indexURL := route.Name2URL("articles.index")
                http.Redirect(w,r,indexURL, http.StatusFound)
            } else {
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "404 文章未找到")
            }
        }
    }
}


