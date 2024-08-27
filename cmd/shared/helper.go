package shared

import (
	"net/http"
)

func GetUniqueIdentifierFromUrl(url string, r *http.Request) *string {
	ref := r.URL.Path[len(url):]
	return &ref
}
