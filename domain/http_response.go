package domain

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	StatusCode int         `json:"status_code"`
	Body       interface{} `json:"body"`
	Errors     interface{} `json:"errors"`
}

func (r *HttpResponse) Write(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(r.StatusCode)
	json.NewEncoder(rw).Encode(r)
}
