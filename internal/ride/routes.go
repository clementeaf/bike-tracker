package ride

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/rides/start", handleStartRide)
	mux.HandleFunc("/rides/end", handleEndRide)
	mux.HandleFunc("/rides/active", handleGetActiveRides)
	mux.HandleFunc("/rides", handleGetAllRides)
	mux.HandleFunc("/rides/", handleGetRideByID)
}
