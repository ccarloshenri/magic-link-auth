package interfaces

import "net/http"

type Controller interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
