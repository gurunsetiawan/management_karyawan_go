package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"karyawan-app/internal/domain"
	"karyawan-app/internal/handler"
)

type MockEmployeeService struct {
	ImportCSVFunc func(csvData io.Reader) (int, []string, error)
}

func (m *MockEmployeeService) GetAllEmployees(page, limit int, search, status string) (*domain.PaginatedEmployeeResponse, error) {
	return nil, nil
}
func (m *MockEmployeeService) GetEmployee(id int) (*domain.Employee, error) { return nil, nil }
func (m *MockEmployeeService) ExportCSV(writer io.Writer) error {
	writer.Write([]byte("Name,Email\nTest,test@test.com\n"))
	return nil
}
func (m *MockEmployeeService) CreateEmployee(employee *domain.Employee) error { return nil }
func (m *MockEmployeeService) UpdateEmployee(employee *domain.Employee) error { return nil }
func (m *MockEmployeeService) DeleteEmployee(id int) error                    { return nil }
func (m *MockEmployeeService) ImportCSV(csvData io.Reader) (int, []string, error) {
	if m.ImportCSVFunc != nil {
		return m.ImportCSVFunc(csvData)
	}
	return 0, nil, nil
}

func TestImportEmployees_Success(t *testing.T) {
	mockSvc := &MockEmployeeService{
		ImportCSVFunc: func(csvData io.Reader) (int, []string, error) {
			return 2, []string{"Baris 3: Error format"}, nil
		},
	}
	h := handler.NewEmployeeHandler(mockSvc)

	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("Name,Email\nBudi,budi@example.com"))
	writer.Close()

	req := httptest.NewRequest("POST", "/employees/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	h.ImportEmployees(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if int(response["success_count"].(float64)) != 2 {
		t.Errorf("Expected success_count 2, got %v", response["success_count"])
	}
	if len(response["failures"].([]interface{})) != 1 {
		t.Errorf("Expected 1 failure, got %v", len(response["failures"].([]interface{})))
	}
}

func TestImportEmployees_NoFile(t *testing.T) {
	mockSvc := &MockEmployeeService{}
	h := handler.NewEmployeeHandler(mockSvc)

	req := httptest.NewRequest("POST", "/employees/import", nil) // no body
	rr := httptest.NewRecorder()

	h.ImportEmployees(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestExportEmployees(t *testing.T) {
	mockSvc := &MockEmployeeService{}
	h := handler.NewEmployeeHandler(mockSvc)

	req, _ := http.NewRequest("GET", "/api/employees/export", nil)
	rr := httptest.NewRecorder()

	h.ExportEmployees(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if contentType := rr.Header().Get("Content-Type"); contentType != "text/csv" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/csv")
	}
	if rr.Body.String() != "Name,Email\nTest,test@test.com\n" {
		t.Errorf("wrong body: %s", rr.Body.String())
	}
}
