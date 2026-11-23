package httpapi

import "net/http"

type healthResponse struct {
	Status string `json:"status"`
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		logRequestRejection(r, "method_not_allowed", "health", 0)
		writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only the GET method is allowed.")
		return
	}

	writeJSONResponse(w, http.StatusOK, healthResponse{Status: "ok"})
}
