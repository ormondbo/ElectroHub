package httpapi

import (
	"log"
	"net/http"

	_ "shop/docs"

	"github.com/swaggo/swag"
)

const swaggerIndex = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Store API documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
        deepLinking: true
      });
    </script>
  </body>
</html>`

func serveSwaggerUI(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/swagger/" && request.URL.Path != "/swagger/index.html" {
		http.NotFound(writer, request)
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := writer.Write([]byte(swaggerIndex)); err != nil {
		log.Printf("write Swagger UI response: %v", err)
	}
}

func serveSwaggerSpec(writer http.ResponseWriter, _ *http.Request) {
	spec, err := swag.ReadDoc()
	if err != nil {
		log.Printf("read generated Swagger specification: %v", err)
		http.Error(writer, "Swagger specification is unavailable", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := writer.Write([]byte(spec)); err != nil {
		log.Printf("write Swagger response: %v", err)
	}
}
