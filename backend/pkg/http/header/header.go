package header

import "net/http"

const (
	ContentType     = "Content-Type"
	JSONContentType = "application/json"
)

func AddJSONContentType(h http.Header) {
	h.Set(ContentType, JSONContentType)
}
