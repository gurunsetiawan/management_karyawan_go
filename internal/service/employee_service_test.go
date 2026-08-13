package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"karyawan-app/internal/domain"
	"karyawan-app/internal/repository"
)

func getTime() time.Time {
	return time.Now()
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"John Doe", "John Doe"},
		{"test@example.com", "test@example.com"},
		{"  test  ", "test"},
		{"", ""},
	}

	for _, test := range tests {
		result := sanitizeInput(test.input)
		if result != test.expected {
			t.Errorf("sanitizeInput(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"", false},
		{"user@domain", false},
	}

	for _, test := range tests {
		result := isValidEmail(test.email)
		if result != test.expected {
			t.Errorf("isValidEmail(%q) = %v, expected %v", test.email, result, test.expected)
		}
	}
}

func TestValidateEmployee(t *testing.T) {
	tests := []struct {
		name        string
		employee    *domain.Employee
		expectError bool
	}{
		{
			name: "valid employee",
			employee: &domain.Employee{
				Name:     "John Doe",
				Email:    "john@example.com",
				Position: "Developer",
				Role:     "Engineer",
				Phone:    "1234567890",
				Alamat:   "123 Main Street, City",
			},
			expectError: false,
		},
		{
			name: "missing name",
			employee: &domain.Employee{
				Email:    "john@example.com",
				Position: "Developer",
				Role:     "Engineer",
				Phone:    "1234567890",
				Alamat:   "123 Main Street, City",
			},
			expectError: true,
		},
		{
			name: "invalid email",
			employee: &domain.Employee{
				Name:     "John Doe",
				Email:    "invalid-email",
				Position: "Developer",
				Role:     "Engineer",
				Phone:    "1234567890",
				Alamat:   "123 Main Street, City",
			},
			expectError: true,
		},
		{
			name: "missing position",
			employee: &domain.Employee{
				Name:     "John Doe",
				Email:    "john@example.com",
				Role:     "Engineer",
				Phone:    "1234567890",
				Alamat:   "123 Main Street, City",
			},
			expectError: true,
		},
		{
			name: "missing address",
			employee: &domain.Employee{
				Name:     "John Doe",
				Email:    "john@example.com",
				Position: "Developer",
				Role:     "Engineer",
				Phone:    "1234567890",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmployee(tt.employee)
			if (err != nil) != tt.expectError {
				t.Errorf("validateEmployee() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestEmployeeStructure(t *testing.T) {
	employee := domain.Employee{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Role:      "Developer",
		Phone:     "1234567890",
		Alamat:    "Test Address",
		CreatedAt: getTime(),
	}

	jsonData, err := json.Marshal(employee)
	if err != nil {
		t.Errorf("Failed to marshal Employee: %v", err)
	}

	var unmarshaled domain.Employee
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal Employee: %v", err)
	}

	if unmarshaled.ID != employee.ID {
		t.Errorf("ID field mismatch: got %d, expected %d", unmarshaled.ID, employee.ID)
	}

	if unmarshaled.Name != employee.Name {
		t.Errorf("Name field mismatch: got %s, expected %s", unmarshaled.Name, employee.Name)
	}

	if unmarshaled.Email != employee.Email {
		t.Errorf("Email field mismatch: got %s, expected %s", unmarshaled.Email, employee.Email)
	}
}

func TestSanitizeEmployee(t *testing.T) {
	employee := &domain.Employee{
		Name:     "  <script>alert('xss')</script>John Doe  ",
		Email:    "  john@example.com  ",
		Position: "<b>Developer</b>",
		Role:     "Engineer",
		Phone:    "1234567890",
		Alamat:   "123 Main St",
	}

	sanitizeEmployee(employee)

	if employee.Name != "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;John Doe" {
		t.Errorf("Name not sanitized properly: got %s", employee.Name)
	}

	if employee.Email != "john@example.com" {
		t.Errorf("Email not sanitized properly: got %s", employee.Email)
	}

	if employee.Position != "&lt;b&gt;Developer&lt;/b&gt;" {
		t.Errorf("Position not sanitized properly: got %s", employee.Position)
	}
}

func TestSanitizeEmployee_EmailPreserved(t *testing.T) {
	employee := &domain.Employee{
		Name:     "John",
		Email:    "o'brien@example.com",
		Position: "Developer",
		Role:     "Engineer",
		Phone:    "1234567890",
		Alamat:   "123 Main St",
	}

	sanitizeEmployee(employee)

	if employee.Email != "o'brien@example.com" {
		t.Errorf("Email was mangled: %s", employee.Email)
	}
}

func TestSanitizeCSVField(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"=CMD|'/C calc'!A0", "'=CMD|'/C calc'!A0"},
		{"+1234", "'+1234"},
		{"-formula", "'-formula"},
		{"@SUM(A1)", "'@SUM(A1)"},
		{"Normal text", "Normal text"},
		{"", ""},
	}

	for _, test := range tests {
		result := repository.SanitizeCSVField(test.input)
		if result != test.expected {
			t.Errorf("SanitizeCSVField(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkSanitizeInput(b *testing.B) {
	input := "<script>alert('xss')</script>"
	for i := 0; i < b.N; i++ {
		sanitizeInput(input)
	}
}

func BenchmarkIsValidEmail(b *testing.B) {
	email := "test@example.com"
	for i := 0; i < b.N; i++ {
		isValidEmail(email)
	}
}

func BenchmarkValidateEmployee(b *testing.B) {
	employee := &domain.Employee{
		Name:     "John Doe",
		Email:    "john@example.com",
		Position: "Developer",
		Role:     "Engineer",
		Phone:    "1234567890",
		Alamat:   "123 Main Street, City",
	}

	for i := 0; i < b.N; i++ {
		_ = validateEmployee(employee)
	}
}

// MockEmployeeRepo is a mock implementation of domain.EmployeeRepository
type MockEmployeeRepo struct {
	CreateFunc func(employee *domain.Employee) error
	FindAllFunc func(page, limit int, search, status string) ([]domain.Employee, int, error)
}

func (m *MockEmployeeRepo) FindAll(page, limit int, search, status string) ([]domain.Employee, int, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(page, limit, search, status)
	}
	return nil, 0, nil
}
func (m *MockEmployeeRepo) FindByID(id int) (*domain.Employee, error) { return nil, nil }
func (m *MockEmployeeRepo) ExportCSV(writer io.Writer) error { return nil }
func (m *MockEmployeeRepo) Create(employee *domain.Employee) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(employee)
	}
	return nil
}
func (m *MockEmployeeRepo) Update(employee *domain.Employee) error { return nil }
func (m *MockEmployeeRepo) Delete(id int) error                    { return nil }

func TestImportCSV_EmptyFile(t *testing.T) {
	mockRepo := &MockEmployeeRepo{}
	svc := NewEmployeeService(mockRepo)

	// Test 1: Empty file should return 0 success, 0 fails, but with an error
	csvData := strings.NewReader("")
	successCount, failures, err := svc.ImportCSV(csvData)

	if err == nil {
		t.Error("Expected error for empty CSV, got nil")
	}
	if successCount != 0 {
		t.Errorf("Expected 0 successes, got %d", successCount)
	}
	if len(failures) != 0 {
		t.Errorf("Expected 0 failures, got %d", len(failures))
	}
}

func TestImportCSV_ValidAndInvalidRows(t *testing.T) {
	var insertedCount int
	mockRepo := &MockEmployeeRepo{
		CreateFunc: func(employee *domain.Employee) error {
			if employee.Email == "duplicate@example.com" {
				return errors.New("email already exists")
			}
			insertedCount++
			return nil
		},
	}
	svc := NewEmployeeService(mockRepo)

	// CSV Header: Name,Email,Position,Role,Phone,Alamat
	csvContent := `Name,Email,Position,Role,Phone,Alamat
Budi,budi@example.com,Manager,Admin,081234,Jakarta
Tono,,Staff,User,082234,Bandung
Andi,duplicate@example.com,Staff,User,083234,Surabaya
`
	csvData := strings.NewReader(csvContent)
	successCount, failures, err := svc.ImportCSV(csvData)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Expect 1 success (Budi)
	if successCount != 1 {
		t.Errorf("Expected 1 success, got %d", successCount)
	}
	if insertedCount != 1 {
		t.Errorf("Expected repo.Create to be called exactly 1 time successfully, got %d", insertedCount)
	}

	// Expect 2 failures (Tono missing email, Andi duplicate email)
	if len(failures) != 2 {
		t.Errorf("Expected 2 failures, got %d", len(failures))
	}
}

func TestExportCSV(t *testing.T) {
	mockRepo := &MockEmployeeRepo{}
	svc := NewEmployeeService(mockRepo)

	var b bytes.Buffer
	err := svc.ExportCSV(&b)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestGetAllEmployees_StatusFilter(t *testing.T) {
	var capturedStatus string
	mockRepo := &MockEmployeeRepo{
		FindAllFunc: func(page, limit int, search, status string) ([]domain.Employee, int, error) {
			capturedStatus = status
			return nil, 0, nil
		},
	}
	svc := NewEmployeeService(mockRepo)

	tests := []string{"active", "inactive", "all", ""}
	for _, status := range tests {
		_, err := svc.GetAllEmployees(1, 10, "", status)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedStatus != status {
			t.Errorf("Expected status %q to be passed to repo, got %q", status, capturedStatus)
		}
	}
}

func TestGetAllEmployees_Pagination(t *testing.T) {
	mockRepo := &MockEmployeeRepo{
		FindAllFunc: func(page, limit int, search, status string) ([]domain.Employee, int, error) {
			return []domain.Employee{{ID: 1}, {ID: 2}}, 12, nil
		},
	}
	svc := NewEmployeeService(mockRepo)
	res, err := svc.GetAllEmployees(2, 5, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Meta.Total != 12 {
		t.Errorf("expected 12 total, got %d", res.Meta.Total)
	}
	if res.Meta.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", res.Meta.TotalPages)
	}
}
