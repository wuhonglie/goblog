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
    "path/filepath"
    "strconv"
    "unicode/utf8"
)

type ArticlesController struct {

}
type ArticlesFormData struct {
    Title, Body string
    URL string
    Errors map[string]string
}
// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request){
    storeURL := route.Name2URL("articles.store")
    data := ArticlesFormData{
        Title: "",
        Body: "",
        URL: storeURL,
        Errors: nil,
    }
    tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
    if err != nil {
        panic(err)
    }
    tmpl.Execute(w, data)
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
        err = tmpl.Execute(w, article)
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
        viewDir := "resources/views"
        files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
        logger.LogError(err)
        newFiles := append(files, viewDir+"/articles/index.gohtml")
        tmpl, err := template.ParseFiles(newFiles...)
        logger.LogError(err)
        err = tmpl.ExecuteTemplate(w, "app", articles)
        logger.LogError(err)

    }
}

func validateArticleFormData(title string, body string) map[string]string {
    errors := make(map[string]string)
    if title == "" {
        errors["title"] = "标题不能为空"
    }else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需介于 3-40"
    }
    if body == "" {
        errors["body"] = "内容不能为空"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度需大于或等于 10 个字节"
    }
    return errors
}
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request){
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")
    errors := validateArticleFormData(title, body)
    if len(errors) == 0 {
        _article := article2.Article{
            Title: title,
            Body: body,
        }
        _article.Create()
        //lastInsertID, err := saveArticleToDB(title, body)
        if _article.ID > 0 {
            fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(_article.ID,10))
        } else {
            //logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "创建文章失败，请联系管理员")
        }
    } else {
        storeURL := route.Name2URL("articles.store")
        data := ArticlesFormData{
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
        }
        tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
        logger.LogError(err)
        tmpl.Execute(w, data)
    }
}

func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
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
        updateURL := route.Name2URL("articles.update", "id",id)
        data := ArticlesFormData{
            Title: article.Title,
            Body: article.Body,
            URL: updateURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
        logger.LogError(err)
        tmpl.Execute(w, data)
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
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")
        errors := validateArticleFormData(title, body)
        if len(errors) == 0 {
            _article.Title = title
            _article.Body = body
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
            updateURL := route.Name2URL("articles.update", "id", id)
            data := ArticlesFormData{
                Title: title,
                Body: body,
                URL: updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
            logger.LogError(err)
            tmpl.Execute(w, data)
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


