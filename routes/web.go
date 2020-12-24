package routes

import (
    "github.com/gorilla/mux"
    "goblog/app/http/controllers"
    "goblog/app/http/middlewares"
    "net/http"
)

func RegisterWebRoutes(r *mux.Router) {
    pc := new(controllers.PagesController)
    // 静态资源
    r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
    r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))
    // 中间件
    //r.Use(middlewares.ForceHTML)
    r.Use(middlewares.StartSession)
    r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
    r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
    // 文章相关页面
    ac := new(controllers.ArticlesController)
    r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
    r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
    r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
    r.HandleFunc("/articles", middlewares.Auth(ac.Store)).Methods("POST").Name("articles.store")
    r.HandleFunc("/articles/create", middlewares.Auth(ac.Create)).Methods("GET").Name("articles.create")
    r.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(ac.Edit)).Methods("GET").Name("articles.edit")
    r.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Update)).Methods("POST").Name("articles.update")
    r.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

    // 用户认证
    uc := new(controllers.UserController)
    r.HandleFunc("/users/{id:[0-9]+}",uc.Show).Methods("GET").Name("users.show")
    auc := new(controllers.AuthController)
    r.HandleFunc("/auth/login", middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
    r.HandleFunc("/auth/dologin", middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
    r.HandleFunc("/auth/logout", middlewares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")
    r.HandleFunc("/auth/register",middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
    r.HandleFunc("/auth/do-register", middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")
}
