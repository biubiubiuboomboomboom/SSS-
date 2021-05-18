package frame

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router{
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string)[]string{
	vs := strings.Split(pattern,"/")

	parts := make([]string,0)

	for _,item := range vs{
		if item != ""{
			parts = append(parts,item)
			if item[0] == '*'{
				break
			}
		}
	}
	return parts
}


func (r *router) addRoute(method string , pattern string , handler HandlerFunc){

	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_,ok := r.roots[method]
	if !ok{
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern,parts,0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string , path string )(*node,map[string]string) {
	parmas := make(map[string]string)
	parts := parsePattern(path)
	root,ok := r.roots[method]
	if !ok{
		return nil,nil
	}
	n := root.search(parts,0)
	if n != nil{
		items := parsePattern(n.pattern)
		for index,item := range items{
			if item[0] == ':'{
				parmas[item[1:]] = parts[index]
			}
			if item[0]=='*' && len(item)>1{
				parmas[item[1:]] = strings.Join(parts[index:],"/")
				break
			}
		}
		return n,parmas
	}
	return nil,nil
}

func (r *router) handle(c *Context)  {
	n, params := r.getRoute(c.method,c.path)
	if n != nil{
		c.params = params
		key := c.method + "-" + n.pattern
		c.handlers = append(c.handlers,r.handlers[key])
	}else {
		c.handlers = append(c.handlers, func(context *Context) {
			context.String(http.StatusNotFound , "404 Not Found : s% \n",context.path)
		})
	}
	c.Next()
}
