package frame

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSON map[string]interface{}

type Context struct {
	// origin
	w http.ResponseWriter
	q *http.Request
	// request info
	path string
	method string
	params map[string]string
	// response info
	statusCode int
	// middleware
	handlers []HandlerFunc
	index int // middleware index
	// engine point can use engine html template
	engine *Engine
}

func newContext(w http.ResponseWriter, q *http.Request) *Context {
	return &Context{
		w: w,
		q: q,
		path: q.URL.Path,
		method: q.Method,
		index: -1,
	}
}

func (c *Context)PostForm(key string)string {
	return c.q.FormValue(key)
}

func (c *Context) Query(key string)string {
	return c.q.URL.Query().Get(key)
}

func (c *Context) Status (code int){
	c.statusCode = code
	c.w.WriteHeader(code)
}

func (c *Context) SetHeader (key string, value string){
	c.w.Header().Set(key,value)
}

func (c *Context) String(code int, format string , values ...interface{}){
	c.SetHeader("Content-Type","text/plain")
	c.Status(code)
	c.w.Write([]byte(fmt.Sprintf(format,values...)))
}

func (c *Context) Json(code int, obj interface{}){
	c.SetHeader("Content-Type","text/json")
	c.Status(code)
	encoder := json.NewEncoder(c.w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.w,err.Error(),500)
	}
}

func (c *Context) Data(code int, data []byte){
	c.Status(code)
	c.w.Write(data)
}

func (c *Context) Html(code int,name string , data interface{}) {
	c.SetHeader("Content-Type","text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.w, name, data); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

func (c *Context) Param(key string) string {
	value ,_ := c.params[key]
	return value
}

func (c *Context) Next(){
	c.index++
	index := len(c.handlers)
	for ;c.index<index ; c.index++{
		c.handlers[c.index](c)
	}
}