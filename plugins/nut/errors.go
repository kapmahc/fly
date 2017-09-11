package nut

import (
	"net/http"
)

// ErrorController handle error
type ErrorController struct {
	ApplicationLayout
}

// Error404 http 404
func (p *ErrorController) Error404() {
	p.show(http.StatusNotFound)
}

// Error500 http 500
func (p *ErrorController) Error500() {
	p.show(http.StatusInternalServerError)
}

func (p *ErrorController) show(c int) {
	p.Data["content"] = http.StatusText(c)
	p.TplName = "error.html"
}
