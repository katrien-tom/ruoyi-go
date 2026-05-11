package route

import "github.com/gin-gonic/gin"

type Meta struct {
	Name       string
	Permission string
	LogTitle   string
}

type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
	Meta        Meta
}

func (r Route) Handlers() []gin.HandlerFunc {
	n := 1 + len(r.Middlewares)
	all := make([]gin.HandlerFunc, 0, n)
	all = append(all, r.Middlewares...)
	all = append(all, r.Handler)
	return all
}

type Group struct {
	Prefix      string
	Middlewares []gin.HandlerFunc
	Routes      []Route
	Children    []*Group
}

func (g Group) Apply(rg *gin.RouterGroup) {
	group := rg
	if g.Prefix != "" {
		group = rg.Group(g.Prefix)
	}
	if len(g.Middlewares) > 0 {
		group.Use(g.Middlewares...)
	}
	for _, r := range g.Routes {
		group.Handle(r.Method, r.Path, r.Handlers()...)
	}
	for _, child := range g.Children {
		child.Apply(group)
	}
}

type Module interface {
	Routes() Group
}

func RegisterModules(rg *gin.RouterGroup, modules ...Module) {
	for _, m := range modules {
		m.Routes().Apply(rg)
	}
}
