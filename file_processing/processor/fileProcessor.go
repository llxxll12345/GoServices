package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func processFile(job *ProcessJob) {
	// Read the content of the uploaded file
	content, err := readFile(job.FilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Modify the content (e.g., convert to uppercase)
	modifiedContent := modifyContent(content)

	// Save the modified content to the processed file
	processedFilePath := "files/processed/" + filepath.Base(job.FilePath)
	err = writeFile(processedFilePath, modifiedContent)
	if err != nil {
		fmt.Println("Error writing processed file:", err)
		return
	}

	mutex.Lock()
	processedFiles[job.FileToken] = processedFilePath
	mutex.Unlock()

	fmt.Println("Processing file:", job.FilePath)
	fmt.Println("File processed:", job.FilePath)
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func modifyContent(content string) string {
	// Modify the content as needed
	// For example, converting to uppercase
	return strings.ToUpper(content)
}

func writeFile(filePath, content string) error {
	time.Sleep(10 * time.Second)
	err := os.WriteFile(filePath, []byte(content), 0644)
	return err
}
