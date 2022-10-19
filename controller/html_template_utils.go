package controller

import "net/http"

type HTMLTemplateExecutor interface {
	ExecuteTemplate(templateName, fileName string, data interface{}, w http.ResponseWriter)
}

type HTMLTemplateExecutorFunc func(templateName, fileName string, data interface{}, w http.ResponseWriter)

func (htef HTMLTemplateExecutorFunc) ExecuteTemplate(templateName, fileName string, data interface{}, w http.ResponseWriter) {
	htef(templateName, fileName, data, w)
}
