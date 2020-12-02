package main

import (
    "database/sql"
    "github.com/gorilla/mux"
    "goblog/bootstrap"
    "goblog/pkg/database"
    "net/http"
    "strings"
    //_ "github.com/go-sql-driver/mysql"
)
var router *mux.Router
var db *sql.DB



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





func main() {
    database.Initialize()
    db = database.DB

    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()
    //router.HandleFunc("/articles/{id:[0-9]+}/delete", articleDeleteHandler).Methods("POST").Name("article.delete")
    router.Use(forceHTMLMiddleware)
    http.ListenAndServe(":3000",  removeTrailingSlash(router))
}
