package errhandler

import (
	"net/http"
	"testing"
	"web"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
		<html>
			<body>
				<h1>Not Found</h1>
			</body>
		</html>
	`))
	builder.AddCode(http.StatusInternalServerError, []byte(`
		<html>
			<body>
				<h1>Internal Server Error</h1>
			</body>
		</html>
	`))

	server := web.NewHttpServer(web.ServerWithMiddleware(builder.Build()))

	server.Start(":8081")
}
