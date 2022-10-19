package external

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileHtmlTemplateUtils = "/external/html_template.go"

func ExecuteHTMLTemplate(templateName, fileName string, data interface{}, w http.ResponseWriter) {
	var callTraceFunc = fmt.Sprintf("%s#ExecuteTemplate", callTraceFileHtmlTemplateUtils)
	tmpl, err := template.New(templateName).ParseFiles(fileName)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrTemplateHTML, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "parsing", fileName))
	}
	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrTemplateHTML, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "execute template", fileName))
	}
}
