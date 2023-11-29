package handlers

import (
	"filemaster/processor"
	"net/http"
)

func Download(w http.ResponseWriter, r *http.Request) {
	fileToken := r.FormValue("token")
	processedFilePath, ok := processor.GetProcessedFilePath(fileToken)
	if ok {
		http.ServeFile(w, r, processedFilePath)
	} else {
		http.Error(w, "File not found or processing not complete", http.StatusNotFound)
	}
}
