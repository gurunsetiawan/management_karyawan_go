-- Buat tabel roles
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Buat tabel users
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

-- Insert role default
INSERT IGNORE INTO roles (name, description) VALUES 
('admin', 'Administrator dengan akses penuh'),
('hrd', 'HRD dengan akses manajemen karyawan'),
('karyawan', 'Karyawan dengan akses terbatas');

-- Insert admin default (password: admin123)
INSERT IGNORE INTO users (id, username, email, password_hash, role_id, is_active, created_at)
VALUES 
(1, 'admin', 'admin@example.com', '$2a$10$Tl0TAod9zFzUyjnuDr2BZu4QNsl39Roq0PLctPg4Enyc/Y4Wtjp6O', 1, 1, NOW()),
(2, 'manager', 'manager@example.com', '$2a$10$Tl0TAod9zFzUyjnuDr2BZu4QNsl39Roq0PLctPg4Enyc/Y4Wtjp6O', 2, 1, NOW()),
(3, 'user', 'user@example.com', '$2a$10$Tl0TAod9zFzUyjnuDr2BZu4QNsl39Roq0PLctPg4Enyc/Y4Wtjp6O', 3, 1, NOW());
