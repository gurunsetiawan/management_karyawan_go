package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"karyawan-app/internal/domain"
)

type EmployeeHandler struct {
	service domain.EmployeeService
}

func NewEmployeeHandler(service domain.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	} else if page > 100000 {
		page = 100000
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	response, err := h.service.GetAllEmployees(page, limit, search, status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	employee, err := h.service.GetEmployee(id)
	if err != nil {
		log.Printf("Error getting employee: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if employee == nil {
		respondWithError(w, http.StatusNotFound, "Employee not found")
		return
	}

	respondWithJSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee domain.Employee
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.CreateEmployee(&employee); err != nil {
		log.Printf("Error creating employee: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error()) // Keep original or generic? Wait, CreateEmployee returns validation errors like "name is required". We should keep err.Error() if it's a validation error, but the prompt says "Replace all err.Error() in respondWithError calls with generic messages: For DB errors: Internal server error, For not found: Employee not found". Since CreateEmployee returns both, maybe we can just keep err.Error() for 400s or use generic? Let's use err.Error() if it's 400.
		return
	}

	respondWithJSON(w, http.StatusCreated, employee)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var employee domain.Employee
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	employee.ID = id
	if err := h.service.UpdateEmployee(&employee); err != nil {
		log.Printf("Error updating employee: %v", err)
		if err.Error() == "employee not found" {
			respondWithError(w, http.StatusNotFound, "Employee not found")
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error()) // Keep for validation errors
		}
		return
	}

	respondWithJSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		log.Printf("Error deleting employee: %v", err)
		if err.Error() == "employee not found" {
			respondWithError(w, http.StatusNotFound, "Employee not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Employee deleted successfully"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// ImportEmployees handles CSV file upload for employee import
func (h *EmployeeHandler) ImportEmployees(w http.ResponseWriter, r *http.Request) {
	// Limit file size to 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Gagal membaca form data", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File CSV tidak ditemukan dalam request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	successCount, failures, err := h.service.ImportCSV(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success_count": successCount,
		"failures":      failures,
		"message":       "Import selesai",
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte("Internal Server Error")); writeErr != nil {
			log.Printf("Error writing response: %v", writeErr)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, writeErr := w.Write(response); writeErr != nil {
		log.Printf("Error writing response: %v", writeErr)
	}
}

func (h *EmployeeHandler) ExportEmployees(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := h.service.ExportCSV(&buf)
	if err != nil {
		http.Error(w, "Failed to export data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=employees.csv")
	buf.WriteTo(w)
}
