package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alexmullins/zip"
)

type ArchiveRequest struct {
	CatalogPath string  `json:"catalog"`
	CatalogTo   *string `json:"save_to"`
	Password    string  `json:"password"`
}

func main() {
	http.HandleFunc("/arch", handleRequest)
	http.ListenAndServe(":8088", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var archiveRequest ArchiveRequest
	err := json.NewDecoder(r.Body).Decode(&archiveRequest)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if archiveRequest.CatalogPath == "" || archiveRequest.Password == "" {
		http.Error(w, "Both catalog and password must be provided", http.StatusBadRequest)
		return
	}

	// Проверяем, что каталог существует
	pathToCatalog := getOSPath(archiveRequest.CatalogTo)

	// Создаем временный файл для сохранения архива
	fileName := fmt.Sprintf("arch_%s.zip", time.Now().Format("20060102150405"))
	tmpFile, err := os.Create(pathToCatalog + fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temporary file: %s", err), http.StatusInternalServerError)
		return
	}
	//defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Создаем архив
	err = createZipArchive(tmpFile, archiveRequest.CatalogPath, archiveRequest.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating archive: %s", err), http.StatusInternalServerError)
		return
	}

	defer func() {
		// Удаляем каталог с файлами
		err = RemoveDirectory(archiveRequest.CatalogPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error removing directory: %s", err), http.StatusInternalServerError)
			return
		}
	}()

	// Отправляем архив в ответ на запрос
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(tmpFile.Name())))
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, tmpFile.Name())
}

func createZipArchive(zipFile *os.File, sourceDir, password string) error {
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err := filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		header := &zip.FileHeader{
			Name:   strings.ReplaceAll(relPath, string(filepath.Separator), "/"),
			Method: zip.Deflate,
			//Modified: info.ModTime(),
		}
		header.SetPassword(password)

		entry, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(entry, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// RemoveDirectory удаляет каталог и его содержимое без подтверждения.
func RemoveDirectory(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.RemoveAll(path)
	})

	if err != nil {
		return err
	}

	return nil
}

func getOSPath(saveTo *string) string {
	if saveTo != nil {
		return *saveTo
	}
	switch runtime.GOOS {
	case "windows":
		return ".\\"
	case "darwin":
		return "./"
	case "linux":
		return "./"
	default:
		return "./"
	}
}
