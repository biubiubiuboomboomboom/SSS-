

This is a simple Web framework prototype.

This program is for learning.

Users can fork to local modifications as needed.


### Quick start 

```go


func main() {
	c := New()
	c.Use(Logger(),Recovery())
	c.GET("/", func(context *Context) {
		context.String(200 , "hello world ! ")
	})
	hello := c.Group("/hello")
	hello.GET("/", func(context *Context) {
			context.String(http.StatusOK, "hello %s, you're at %s \n", context.Query("name"), context.path)
		})

	hello.GET("/:name", func(context *Context) {
			context.String(http.StatusOK, "hello %s \n", context.Param("name"))
		})
	c.Start(":9000")
}


```