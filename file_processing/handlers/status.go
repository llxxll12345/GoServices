package handlers

import (
	"filemaster/processor"
	"fmt"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	fileToken := r.FormValue("token")
	_, ok := processor.GetProcessedFilePath(fileToken)
	if ok {
		fmt.Fprintln(w, "File processing complete. You can download your file.")
	} else {
		fmt.Fprintln(w, "File processing is still in progress. Please check again later.")
	}
}
