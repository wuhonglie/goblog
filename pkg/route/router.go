package route

import (
    "github.com/gorilla/mux"
    "goblog/pkg/logger"
    "net/http"
)

var route *mux.Router

func SetRoute(r *mux.Router) {
    route = r
}
func Name2URL(routeName string, pairs ...string) string {
    url, err := route.Get(routeName).URL(pairs...)
    if err != nil {
        logger.LogError(err)
        return ""
    }

    return url.String()
}
//GetRouteVariable 获取 URI 路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}
