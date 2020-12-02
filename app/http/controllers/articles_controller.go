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


