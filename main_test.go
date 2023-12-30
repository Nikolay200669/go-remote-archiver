package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	// Подготовка тестовых данных
	testData := ArchiveRequest{
		CatalogPath: "test",
		Password:    "test_password",
	}
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}

	// Создание виртуального HTTP-сервера
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	// Выполнение запроса
	handler.ServeHTTP(rr, req)

	// Проверка кода состояния
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Дополнительные проверки могут быть добавлены в соответствии с вашими требованиями
}

//func TestCreateZipArchive(t *testing.T) {
//	// Создание виртуального файла
//	tmpFile, err := createTempFile("test_archive*.zip")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer os.Remove(tmpFile.Name())
//	defer tmpFile.Close()
//
//	// Подготовка тестовых данных
//	testData := ArchiveRequest{
//		CatalogPath: "test_catalog",
//		Password:    "test_password",
//	}
//
//	// Вызов функции createZipArchive
//	err = createZipArchive(tmpFile, testData.CatalogPath, testData.Password)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Дополнительные проверки могут быть добавлены в соответствии с вашими требованиями
//}
//
//func createTempFile(pattern string) (*os.File, error) {
//	tmpFile, err := os.CreateTemp("", pattern)
//	if err != nil {
//		return nil, err
//	}
//	return tmpFile, nil
//}
