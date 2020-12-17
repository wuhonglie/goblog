package view

import (
    "goblog/pkg/logger"
    "goblog/pkg/route"
    "html/template"
    "io"
    "path/filepath"
    "strings"
)
type D map[string]interface{}
// Render 渲染通用视图
func Render(w io.Writer, data interface{}, tplFiles ...string) {
    RenderTemplate(w, "app", data, tplFiles...)
}
// RenderSimple 渲染简单的视图
func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
    RenderTemplate(w, "simple", data, tplFiles...)
}

func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
    viewDir := "resources/views/"
    for i,f := range tplFiles {
        tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
    }
    layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
    logger.LogError(err)
    allFiles := append(layoutFiles, tplFiles...)
    tmpl, err := template.New("").Funcs(template.FuncMap{
        "RouteName2URL" : route.Name2URL,
    }).ParseFiles(allFiles...)
    tmpl.ExecuteTemplate(w, name, data)
}
