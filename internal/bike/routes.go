package bike

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/bikes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			HandleGetAllBikes(w, r)
		case http.MethodPost:
			HandleRegisterBike(w, r)
		default:
			http.Error(w, `{"error": "Método no permitido"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/bikes/available", HandleGetAvailableBikes)
	mux.HandleFunc("/bikes/status", HandleUpdateBikeStatus)
}
