package rest

import (
	"encoding/json"
	"net/http"
)

func writeOkResponse(w http.ResponseWriter, data interface{}) {
	jsonData, _ := json.Marshal(data)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getHealthResponseData() map[string]string {
	data := make(map[string]string)
	data["status"] = "UP"

	return data
}
