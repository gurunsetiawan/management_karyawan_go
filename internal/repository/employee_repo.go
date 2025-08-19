package repository

import (
	"database/sql"
	"time"

	"karyawan-app/internal/domain"
)

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) FindAll() ([]domain.Employee, error) {
	query := `SELECT id, name, email, role, phone, alamat, created_at FROM employees ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var e domain.Employee
		var createdAtStr string
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Phone, &e.Alamat, &createdAtStr); err != nil {
			return nil, err
		}
		// Convert string time to time.Time if needed
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *employeeRepository) FindByID(id int) (*domain.Employee, error) {
	query := `SELECT id, name, email, role, phone, alamat, created_at FROM employees WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var e domain.Employee
	var createdAtStr string
	err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Phone, &e.Alamat, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	return &e, nil
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	query := `INSERT INTO employees (name, email, role, phone, alamat) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, employee.Name, employee.Email, employee.Role, employee.Phone, employee.Alamat)
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
	query := `UPDATE employees SET name=?, email=?, role=?, phone=?, alamat=? WHERE id=?`
	_, err := r.db.Exec(query, employee.Name, employee.Email, employee.Role, employee.Phone, employee.Alamat, employee.ID)
	return err
}

func (r *employeeRepository) Delete(id int) error {
	query := `DELETE FROM employees WHERE id=?`
	_, err := r.db.Exec(query, id)
	return err
}
