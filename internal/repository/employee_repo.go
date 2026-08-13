package repository

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"karyawan-app/internal/domain"
)

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

// escapeLike removes wildcard characters for exact match in LIKE
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

// SanitizeCSVField prevents CSV injection
func SanitizeCSVField(s string) string {
	if s == "" {
		return s
	}
	if strings.HasPrefix(s, "=") || strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") || strings.HasPrefix(s, "@") || strings.HasPrefix(s, "\t") || strings.HasPrefix(s, "\r") {
		return "'" + s
	}
	return s
}

func (r *employeeRepository) FindAll(page, limit int, search, status string) ([]domain.Employee, int, error) {
	// First get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM employees WHERE 1=1`
	var countArgs []interface{}

	if status == "active" || status == "" {
		countQuery += ` AND deleted_at IS NULL`
	} else if status == "inactive" {
		countQuery += ` AND deleted_at IS NOT NULL`
	}
	
	if search != "" {
		countQuery += ` AND (name LIKE ? OR email LIKE ?)`
		searchParam := "%" + escapeLike(search) + "%"
		countArgs = append(countArgs, searchParam, searchParam)
	}
	
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	query := `SELECT id, name, email, position, role, phone, alamat, created_at, updated_at FROM employees WHERE 1=1`
	var args []interface{}

	if status == "active" || status == "" {
		query += ` AND deleted_at IS NULL`
	} else if status == "inactive" {
		query += ` AND deleted_at IS NOT NULL`
	}
	
	if search != "" {
		query += ` AND (name LIKE ? OR email LIKE ?)`
		searchParam := "%" + escapeLike(search) + "%"
		args = append(args, searchParam, searchParam)
	}
	
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var e domain.Employee
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Position, &e.Role, &e.Phone, &e.Alamat, &createdAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if createdAt.Valid {
			e.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			e.UpdatedAt = updatedAt.Time
		}
		employees = append(employees, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	
	if employees == nil {
		employees = []domain.Employee{}
	}
	
	return employees, total, nil
}

func (r *employeeRepository) FindByID(id int) (*domain.Employee, error) {
	query := `SELECT id, name, email, position, role, phone, alamat, created_at, updated_at FROM employees WHERE id = ? AND deleted_at IS NULL`
	row := r.db.QueryRow(query, id)

	var e domain.Employee
	var createdAt, updatedAt sql.NullTime
	err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Position, &e.Role, &e.Phone, &e.Alamat, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if createdAt.Valid {
		e.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		e.UpdatedAt = updatedAt.Time
	}
	return &e, nil
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	query := `INSERT INTO employees (name, email, position, role, phone, alamat) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, employee.Name, employee.Email, employee.Position, employee.Role, employee.Phone, employee.Alamat)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	employee.ID = int(id)
	return nil
}

func (r *employeeRepository) Update(employee *domain.Employee) error {
	query := `UPDATE employees SET name=?, email=?, position=?, role=?, phone=?, alamat=?, updated_at=NOW() WHERE id=? AND deleted_at IS NULL`
	result, err := r.db.Exec(query, employee.Name, employee.Email, employee.Position, employee.Role, employee.Phone, employee.Alamat, employee.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *employeeRepository) Delete(id int) error {
	// Soft delete - set deleted_at timestamp instead of actually deleting
	query := `UPDATE employees SET deleted_at=NOW() WHERE id=? AND deleted_at IS NULL`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *employeeRepository) ExportCSV(writer io.Writer) error {
	query := "SELECT name, email, position, role, phone, alamat FROM employees WHERE deleted_at IS NULL ORDER BY id ASC"
	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	csvWriter := csv.NewWriter(writer)
	// Don't defer csvWriter.Flush(), just call it before returning
	
	if err := csvWriter.Write([]string{"Name", "Email", "Position", "Role", "Phone", "Alamat"}); err != nil {
		return err
	}

	for rows.Next() {
		var name, email, position, role, phone, alamat string
		if err := rows.Scan(&name, &email, &position, &role, &phone, &alamat); err != nil {
			return err
		}
		
		name = SanitizeCSVField(name)
		email = SanitizeCSVField(email)
		position = SanitizeCSVField(position)
		role = SanitizeCSVField(role)
		phone = SanitizeCSVField(phone)
		alamat = SanitizeCSVField(alamat)

		if err := csvWriter.Write([]string{name, email, position, role, phone, alamat}); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	
	csvWriter.Flush()
	return csvWriter.Error()
}
