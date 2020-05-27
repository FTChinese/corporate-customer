package controller

import "net/http"

func RenderIndex(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(B2BIndexHTML))
	if err != nil {
		return err
	}

	return nil
}
