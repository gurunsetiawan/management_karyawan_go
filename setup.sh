#!/bin/bash

echo "========================================"
echo "KaryawanApp - Database Setup Script"
echo "========================================"
echo ""

# Get MySQL password from user
read -p "Enter MySQL root password (press Enter if no password): " MYSQL_PASSWORD

# Set default password if empty
if [ -z "$MYSQL_PASSWORD" ]; then
    MYSQL_PASSWORD=""
    echo "Using no password for MySQL root..."
else
    echo "Using provided password..."
fi

export MYSQL_PWD="$MYSQL_PASSWORD"

echo ""
echo "Creating database..."
mysql -u root -e "CREATE DATABASE IF NOT EXISTS karyawan_db;" 2>&1

if [ $? -eq 0 ]; then
    echo "✓ Database 'karyawan_db' created successfully!"
else
    echo "✗ Failed to create database. Check MySQL credentials."
    exit 1
fi

echo ""
echo "Updating .env file..."
# Update .env with MySQL password
cat > .env << EOF
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=$MYSQL_PASSWORD
DB_NAME=karyawan_db
PORT=8083
HOST=127.0.0.1
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
JWT_SECRET=$(openssl rand -base64 32)
JWT_EXPIRATION_HOURS=24
EOF

echo "✓ .env file updated!"

echo ""
echo "Building application..."
go build -o karyawan-app cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
else
    echo "✗ Build failed!"
    exit 1
fi

echo ""
echo "Running seed command..."
go run cmd/seed/seed.go

if [ $? -eq 0 ]; then
    echo "✓ Seed command executed!"
else
    echo "⚠ Seed command failed. You can run it manually later."
fi

echo ""
echo "========================================"
echo "Setup Complete!"
echo "========================================"
echo ""
echo "Default Login Credentials:"
echo "  Username: admin"
echo "  Password: admin123"
echo ""
echo "To start the application:"
echo "  ./karyawan-app"
echo "  or"
echo "  make run"
echo ""
echo "Access the application at:"
echo "  http://127.0.0.1:8083"
echo "========================================"
