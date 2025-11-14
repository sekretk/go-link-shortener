package middleware

import "net/http"

type WrapperWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrapperWriter) WriteHandler(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}
