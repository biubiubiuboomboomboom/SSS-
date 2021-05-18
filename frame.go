package frame

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
	// html render
	htmlTemplates *template.Template
	funcMap template.FuncMap
}

type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup // nesting
	engine *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// -----------------------
// engine function
// -----------------------

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandlerFunc
	for _,group := range engine.groups{
		if strings.HasPrefix(request.URL.Path,group.prefix){
			middlewares = append(middlewares,group.middlewares...)
		}
	}
	c := newContext(writer,request)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

func (engine *Engine)addRoute (method string , pattern string , handler HandlerFunc){
	engine.router.addRoute(method,pattern,handler)
}

func (engine *Engine)GET(pattern string , handler HandlerFunc){
	engine.addRoute("GET",pattern,handler)
}

func (engine *Engine)POST(pattern string , handler HandlerFunc){
	engine.addRoute("POST",pattern,handler)
}

func (engine *Engine)Start(addr string)(err error){
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// -----------------------
// gourp  function
// -----------------------

func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix+prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups,newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute (method string , str string , handler HandlerFunc){
	log.Printf("Route %4s - %s", method, group.prefix+str)
	group.engine.router.addRoute(method,group.prefix+str,handler)
}

func (group *RouterGroup) GET(pattern string , handler HandlerFunc){
	group.addRoute("GET",pattern,handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares,middlewares...)
}

func (group *RouterGroup) createStaticHandler(relativePath string,fs http.FileSystem) HandlerFunc{
	absolutePath := path.Join(group.prefix,relativePath)
	fileServer := http.StripPrefix(absolutePath,http.FileServer(fs))
	return func(context *Context) {
		file := context.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			context.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(context.w,context.q)
	}
}

func (group *RouterGroup) Static(relativePath string , root string){
	handler := group.createStaticHandler(relativePath,http.Dir(root))
	urlPattern := path.Join(relativePath,"/*filepath")
	group.GET(urlPattern,handler)
}