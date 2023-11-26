package router

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc `json:"-"`
}
type Routes map[string]Route

func (r Route) UpdateHandlerFunc(f http.HandlerFunc) Route {
	r.HandlerFunc = f
	return r
}

func (r *Routes) UpdateHandlerFunc(key string, f http.HandlerFunc) {
	(*r)[key] = (*r)[key].UpdateHandlerFunc(f)
}
