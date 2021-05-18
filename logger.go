package frame

import (
	"log"
	"time"
)

func Logger() HandlerFunc{
	return func(context *Context) {
		t:= time.Now()
		context.Next()
		log.Printf("[%d] %s in %v", context.statusCode, context.q.RequestURI, time.Since(t))
	}
}
