package main

import (
    "database/sql"
    "fmt"
    "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "goblog/pkg/logger"
    "goblog/pkg/route"
    "goblog/pkg/types"
    "html/template"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
    "unicode/utf8"
    //_ "github.com/go-sql-driver/mysql"
)
var router *mux.Router
var db *sql.DB
func initDB() {
    var err error
    config := mysql.Config{
        User: "homestead",
        Passwd: "secret",
        Addr: "127.0.0.1:33060",
        Net:    "tcp",
        DBName: "goblog",
        AllowNativePasswords: true,
    }
    db, err = sql.Open("mysql", config.FormatDSN())
    logger.LogError(err)
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    err = db.Ping()
    logger.LogError(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<h1>Hello，欢迎来到 goblog! </h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 " +
        "<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

type Article struct {
    Title, Body string
    ID int64
}
func (a Article) Link() string {
    showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID,10))
    if err != nil {
        logger.LogError(err)
        return ""
    }
    return showURL.String()
}
func (a Article) Delete() (rowsAffected int64, err error) {
    rs, err := db.Exec("DELETE FROM articles WHERE id = "+strconv.FormatInt(a.ID, 10))
    if err != nil {
        return 0, err
    }
    if n, _ := rs.RowsAffected(); n > 0 {
        return n, nil
    }
    return 0, nil
}
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    article, err := getArticleByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
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

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM articles")
    logger.LogError(err)
    defer rows.Close()
    var articles []Article
    for rows.Next() {
        var article Article
        err := rows.Scan(&article.ID, &article.Title, &article.Body)
        logger.LogError(err)
        articles = append(articles, article)
    }
    err = rows.Err()
    logger.LogError(err)
    tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
    logger.LogError(err)
    tmpl.Execute(w, articles)
}
type ArticlesFormData struct {
    Title, Body string
    URL *url.URL
    Errors map[string]string
}
func articleStoreHandler(w http.ResponseWriter, r *http.Request) {
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")
    errors := make(map[string]string)
    if title == "" {
        errors["title"] = "标题不能为空"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需介于 3-40"
    }
    if body == "" {
        errors["body"] = "内容不能为空"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度需大于或等于10个字节"
    }
    if len(errors) == 0 {
        lastInsertID, err := saveArticleToDB(title, body)
        if lastInsertID > 0 {
            fmt.Fprintf(w, "插入成功,ID为"+strconv.FormatInt(lastInsertID,10))
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        storeURL, _ := router.Get("articles.store").URL()
        data := ArticlesFormData{
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
        }
        tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
        if err != nil {
            panic(err)
        }
        tmpl.Execute(w, data)
    }
}

func saveArticleToDB(title string, body string) (int64, error) {
    var (
        id int64
        err error
        rs sql.Result
        stmt *sql.Stmt
    )
    stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()
    rs, err = stmt.Exec(title, body)
    if err != nil {
        return 0, err
    }
    if id, err = rs.LastInsertId(); id > 0 {
        return id, nil
    }
    return 0, err
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Content-Type", "text/html;charset=utf-8")
        next.ServeHTTP(w,r)
    })
}

func removeTrailingSlash(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
        }
        next.ServeHTTP(w, r)
    })
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

    storeURL, _ := router.Get("articles.store").URL()
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

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    article,err := getArticleByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        updateURL, _ := router.Get("articles.update").URL("id", id)
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

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request){
     id := route.GetRouteVariable("id", r)
     _, err := getArticleByID(id)
     if err != nil {
         if err == sql.ErrNoRows {
             w.WriteHeader(http.StatusNotFound)
             fmt.Fprint(w, "404 文章未找到")
         }else{
             logger.LogError(err)
             w.WriteHeader(http.StatusInternalServerError)
             fmt.Fprint(w, "500 服务器内部错误")
         }
     }else{
         title := r.PostFormValue("title")
         body := r.PostFormValue("body")
         errors := make(map[string]string)
         if title == "" {
             errors["title"] = "标题不能为空"
         } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
             errors["title"] = "标题长度需介于3-40"
         }
         if body == "" {
             errors["body"] = "内容不能为空"
         }else if utf8.RuneCountInString(body) < 10 {
             errors["body"] = "内容长度需大于或等于10个字符"
         }
         if len(errors) == 0 {
             query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
             rs, err := db.Exec(query, title, body, id)
             if err != nil {
                 logger.LogError(err)
                 w.WriteHeader(http.StatusInternalServerError)
                 fmt.Fprint(w, "500 服务器内部错误")
             }
             if n, _ := rs.RowsAffected(); n > 0 {
                 showURL, _ := router.Get("articles.show").URL("id",id)
                 http.Redirect(w, r, showURL.String(),http.StatusFound)
             } else{
                 fmt.Fprint(w, "您没有做任何更改")
             }
         } else {
             updateURL,_ := router.Get("articles.edit").URL("id",id)
             data := ArticlesFormData{
                 Title: title,
                 Body: body,
                 URL: updateURL,
                 Errors: errors,
             }
             tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
             logger.LogError(err)
             tmpl.Execute(w,data)
         }
     }
}

func articleDeleteHandler(w http.ResponseWriter, r *http.Request){
    id := route.GetRouteVariable("id", r)
    article, err := getArticleByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        }else{
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    }else {
        rowsAffected, err := article.Delete()
        if err != nil {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        } else {
            if rowsAffected > 0 {
                indexURL, _ := router.Get("articles.index").URL()
                http.Redirect(w, r, indexURL.String(), http.StatusFound)
            } else {
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "404 文章未找到")
            }
        }
    }
}

func createTables() {
    createArticlesSQL := `CREATE TABLE IF NOT EXISTS
articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
);
`
    _, err := db.Exec(createArticlesSQL)
    logger.LogError(err)
}

func getArticleByID(id string) (Article, error) {
    article := Article{}
    query := "SELECT * FROM articles WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
    return article, err
}


func main() {
    initDB()
    createTables()
    route.Initialize()
    router = route.Router
    router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articleStoreHandler).Methods("POST").Name("articles.store")
    router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
    router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
    router.HandleFunc("/articles/{id:[0-9]+}/delete", articleDeleteHandler).Methods("POST").Name("article.delete")
    // 自定义 404 页面
    router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
    router.Use(forceHTMLMiddleware)
    http.ListenAndServe(":3000",  removeTrailingSlash(router))
}
