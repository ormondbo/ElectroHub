package swagger

import (
	_ "electrohub/internal/docs"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

func ServeSwaggerUI() http.Handler {
	return httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
}

func ServeSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	doc, err := swag.ReadDoc()
	if err != nil {
		http.Error(w, "failed to read swagger documentation", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(doc))
}
