-- Add position column to employees table
ALTER TABLE employees
ADD COLUMN position VARCHAR(100) NOT NULL AFTER email;

-- Update existing records with default position based on role
UPDATE employees SET position = 
  CASE 
    WHEN role = 'Manager' THEN 'Senior Manager'
    WHEN role = 'Developer' THEN 'Software Engineer'
    WHEN role = 'Designer' THEN 'UI/UX Designer'
    WHEN role = 'HR' THEN 'HR Specialist'
    WHEN role = 'QA' THEN 'Quality Assurance Engineer'
    ELSE 'Staff'
  END;

-- Make sure the column is not null after updating all records
ALTER TABLE employees MODIFY COLUMN position VARCHAR(100) NOT NULL;
