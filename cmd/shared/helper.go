package shared

import (
	"net/http"
)

func GetPartReferenceFromUrl(url string, r *http.Request) *string {
	ref := r.URL.Path[len(url):]
	return &ref
}
