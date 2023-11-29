package handlers

import (
	"filemaster/processor"
	"filemaster/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileToken := utils.GenerateToken(handler.Filename)

	filePath := "files/uploaded/" + fileToken + filepath.Ext(handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	processor.AddJob(&processor.ProcessJob{FileToken: fileToken, FilePath: filePath})

	fmt.Fprintln(w, "File uploaded successfully. Processing request. Token:", fileToken)
}
