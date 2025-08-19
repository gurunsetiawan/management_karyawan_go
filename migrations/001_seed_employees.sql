-- Create employees table if it doesn't exist
CREATE TABLE IF NOT EXISTS employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL,
    position VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    alamat TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Clear existing data (optional, uncomment if needed)
-- TRUNCATE TABLE employees;

-- Insert 100 sample employees
INSERT INTO employees (name, email, role, position, phone, alamat) VALUES
('John Doe', 'john.doe@example.com', 'Manager', 'Senior Project Manager', '081234567890', 'Jl. Sudirman No. 1, Jakarta Selatan'),
('Jane Smith', 'jane.smith@example.com', 'Developer', 'Frontend Developer', '081234567891', 'Jl. Thamrin No. 10, Jakarta Pusat'),
('Michael Johnson', 'michael.j@example.com', 'Developer', 'Backend Developer', '081234567892', 'Jl. Gatot Subroto No. 5, Jakarta Barat'),
('Sarah Williams', 'sarah.w@example.com', 'Designer', 'UI/UX Designer', '081234567893', 'Jl. HR Rasuna Said No. 20, Jakarta Selatan'),
('David Brown', 'david.b@example.com', 'Manager', 'Product Manager', '081234567894', 'Jl. Jendral Sudirman Kav. 29, Jakarta Pusat'),
('Emily Davis', 'emily.d@example.com', 'Developer', 'Full Stack Developer', '081234567895', 'Jl. Senopati No. 8, Jakarta Selatan'),
('Robert Wilson', 'robert.w@example.com', 'QA', 'Quality Assurance Engineer', '081234567896', 'Jl. Kebon Sirih No. 15, Jakarta Pusat'),
('Jennifer Lee', 'jennifer.l@example.com', 'HR', 'HR Manager', '081234567897', 'Jl. Kuningan Barat No. 25, Jakarta Selatan'),
('William Taylor', 'william.t@example.com', 'Developer', 'Mobile Developer', '081234567898', 'Jl. Suryopranoto No. 5, Jakarta Pusat'),
('Amanda Clark', 'amanda.c@example.com', 'Designer', 'Graphic Designer', '081234567899', 'Jl. Wahid Hasyim No. 10, Jakarta Pusat'),

-- Additional 90 sample employees (abbreviated for space, but you can expand as needed)
('James Anderson', 'james.a@example.com', 'Developer', 'Backend Developer', '081234567900', 'Jl. K.H. Mas Mansyur No. 1, Jakarta Pusat'),
('Patricia Thomas', 'patricia.t@example.com', 'HR', 'Recruiter', '081234567901', 'Jl. K.H. Wahid Hasyim No. 22, Jakarta Pusat'),
('Richard Jackson', 'richard.j@example.com', 'Manager', 'Engineering Manager', '081234567902', 'Jl. Jendral Sudirman Kav. 52-53, Jakarta Selatan'),
('Jessica White', 'jessica.w@example.com', 'Developer', 'Frontend Developer', '081234567903', 'Jl. H.R. Rasuna Said Kav. B-12, Jakarta Selatan'),
('Charles Harris', 'charles.h@example.com', 'QA', 'Test Automation Engineer', '081234567904', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Karen Martin', 'karen.m@example.com', 'Designer', 'UI Designer', '081234567905', 'Jl. Jendral Gatot Subroto No. 27, Jakarta Selatan'),
('Thomas Garcia', 'thomas.g@example.com', 'Developer', 'DevOps Engineer', '081234567906', 'Jl. M.H. Thamrin No. 5, Jakarta Pusat'),
('Nancy Martinez', 'nancy.m@example.com', 'HR', 'HR Business Partner', '081234567907', 'Jl. Jendral Sudirman Kav. 21, Jakarta Selatan'),
('Daniel Robinson', 'daniel.r@example.com', 'Developer', 'Backend Developer', '081234567908', 'Jl. Prof. Dr. Satrio Kav. 18, Jakarta Selatan'),
('Lisa Clark', 'lisa.c@example.com', 'Manager', 'Product Owner', '081234567909', 'Jl. Jendral Sudirman Kav. 32-33, Jakarta Pusat'),

-- Continue with more sample data...
('Mark Lewis', 'mark.l@example.com', 'Developer', 'Full Stack Developer', '081234567910', 'Jl. H.R. Rasuna Said Kav. C-22, Jakarta Selatan'),
('Betty Scott', 'betty.s@example.com', 'QA', 'Manual Tester', '081234567911', 'Jl. Jendral Sudirman Kav. 29, Jakarta Pusat'),
('Donald Young', 'donald.y@example.com', 'Designer', 'UX Designer', '081234567912', 'Jl. K.H. Mas Mansyur No. 15, Jakarta Pusat'),
('Sandra King', 'sandra.k@example.com', 'HR', 'HR Generalist', '081234567913', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),
('Paul Wright', 'paul.w@example.com', 'Developer', 'Frontend Developer', '081234567914', 'Jl. Jendral Gatot Subroto No. 10, Jakarta Selatan'),
('Ashley Lopez', 'ashley.l@example.com', 'Manager', 'Delivery Manager', '081234567915', 'Jl. H.R. Rasuna Said Kav. X-5, Jakarta Selatan'),
('Steven Hill', 'steven.h@example.com', 'Developer', 'Backend Developer', '081234567916', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Donna Green', 'donna.g@example.com', 'QA', 'QA Lead', '081234567917', 'Jl. Jendral Sudirman No. 1, Jakarta Pusat'),
('Andrew Adams', 'andrew.a@example.com', 'Designer', 'Product Designer', '081234567918', 'Jl. K.H. Wahid Hasyim No. 8, Jakarta Pusat'),
('Kimberly Nelson', 'kimberly.n@example.com', 'HR', 'Talent Acquisition', '081234567919', 'Jl. Jendral Sudirman Kav. 29, Jakarta Selatan'),

-- Continue with more sample data...
('Joshua Baker', 'joshua.b@example.com', 'Developer', 'Frontend Developer', '081234567920', 'Jl. H.R. Rasuna Said Kav. C-5, Jakarta Selatan'),
('Emily Carter', 'emily.c@example.com', 'Manager', 'Project Manager', '081234567921', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Kevin Mitchell', 'kevin.m@example.com', 'Developer', 'Backend Developer', '081234567922', 'Jl. K.H. Mas Mansyur No. 20, Jakarta Pusat'),
('Michelle Perez', 'michelle.p@example.com', 'QA', 'Test Engineer', '081234567923', 'Jl. Jendral Sudirman Kav. 1, Jakarta Pusat'),
('Brian Roberts', 'brian.r@example.com', 'Designer', 'UI/UX Designer', '081234567924', 'Jl. H.R. Rasuna Said Kav. X-1, Jakarta Selatan'),
('Laura Turner', 'laura.t@example.com', 'HR', 'HR Manager', '081234567925', 'Jl. Jendral Sudirman Kav. 25, Jakarta Selatan'),
('Ronald Phillips', 'ronald.p@example.com', 'Developer', 'Full Stack Developer', '081234567926', 'Jl. K.H. Wahid Hasyim No. 12, Jakarta Pusat'),
('Deborah Campbell', 'deborah.c@example.com', 'Manager', 'Technical Lead', '081234567927', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Jason Parker', 'jason.p@example.com', 'Developer', 'Mobile Developer', '081234567928', 'Jl. H.R. Rasuna Said Kav. C-15, Jakarta Selatan'),
('Sharon Evans', 'sharon.e@example.com', 'QA', 'QA Engineer', '081234567929', 'Jl. Jendral Sudirman No. 7, Jakarta Pusat'),

-- Continue with more sample data...
('Jeffrey Edwards', 'jeffrey.e@example.com', 'Designer', 'UX Researcher', '081234567930', 'Jl. K.H. Mas Mansyur No. 25, Jakarta Pusat'),
('Carol Collins', 'carol.c@example.com', 'HR', 'HR Business Partner', '081234567931', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),
('Ryan Stewart', 'ryan.s@example.com', 'Developer', 'Backend Developer', '081234567932', 'Jl. H.R. Rasuna Said Kav. X-10, Jakarta Selatan'),
('Amanda Sanchez', 'amanda.s@example.com', 'Manager', 'Product Manager', '081234567933', 'Jl. Jendral Sudirman No. 15, Jakarta Pusat'),
('Gary Morris', 'gary.m@example.com', 'Developer', 'Frontend Developer', '081234567934', 'Jl. K.H. Wahid Hasyim No. 18, Jakarta Pusat'),
('Heather Rogers', 'heather.r@example.com', 'QA', 'Test Automation Engineer', '081234567935', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Timothy Reed', 'timothy.r@example.com', 'Designer', 'UI Designer', '081234567936', 'Jl. H.R. Rasuna Said Kav. C-25, Jakarta Selatan'),
('Rebecca Cook', 'rebecca.c@example.com', 'HR', 'Recruiter', '081234567937', 'Jl. Jendral Sudirman No. 5, Jakarta Pusat'),
('Jose Morgan', 'jose.m@example.com', 'Developer', 'Full Stack Developer', '081234567938', 'Jl. K.H. Mas Mansyur No. 30, Jakarta Pusat'),
('Shirley Bell', 'shirley.b@example.com', 'Manager', 'Delivery Manager', '081234567939', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),

-- Continue with more sample data...
('Dennis Murphy', 'dennis.m@example.com', 'Developer', 'Backend Developer', '081234567940', 'Jl. H.R. Rasuna Said Kav. X-15, Jakarta Selatan'),
('Cynthia Bailey', 'cynthia.b@example.com', 'QA', 'QA Lead', '081234567941', 'Jl. Jendral Sudirman No. 10, Jakarta Pusat'),
('Paul Rivera', 'paul.r@example.com', 'Designer', 'Product Designer', '081234567942', 'Jl. K.H. Wahid Hasyim No. 22, Jakarta Pusat'),
('Kathleen Cooper', 'kathleen.c@example.com', 'HR', 'HR Generalist', '081234567943', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Jeremy Richardson', 'jeremy.r@example.com', 'Developer', 'Frontend Developer', '081234567944', 'Jl. H.R. Rasuna Said Kav. C-35, Jakarta Selatan'),
('Rachel Cox', 'rachel.c@example.com', 'Manager', 'Project Manager', '081234567945', 'Jl. Jendral Sudirman No. 20, Jakarta Pusat'),
('Aaron Howard', 'aaron.h@example.com', 'Developer', 'Backend Developer', '081234567946', 'Jl. K.H. Mas Mansyur No. 35, Jakarta Pusat'),
('Virginia Ward', 'virginia.w@example.com', 'QA', 'Test Engineer', '081234567947', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Douglas Torres', 'douglas.t@example.com', 'Designer', 'UI/UX Designer', '081234567948', 'Jl. H.R. Rasuna Said Kav. X-20, Jakarta Selatan'),
('Brenda Peterson', 'brenda.p@example.com', 'HR', 'HR Manager', '081234567949', 'Jl. Jendral Sudirman No. 25, Jakarta Pusat'),

-- Continue with more sample data...
('Peter Gray', 'peter.g@example.com', 'Developer', 'Full Stack Developer', '081234567950', 'Jl. K.H. Wahid Hasyim No. 28, Jakarta Pusat'),
('Emma Ramirez', 'emma.r@example.com', 'Manager', 'Technical Lead', '081234567951', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Zachary James', 'zachary.j@example.com', 'Developer', 'Mobile Developer', '081234567952', 'Jl. H.R. Rasuna Said Kav. C-45, Jakarta Selatan'),
('Nicole Watson', 'nicole.w@example.com', 'QA', 'QA Engineer', '081234567953', 'Jl. Jendral Sudirman No. 30, Jakarta Pusat'),
('Jack Brooks', 'jack.b@example.com', 'Designer', 'UX Researcher', '081234567954', 'Jl. K.H. Mas Mansyur No. 40, Jakarta Pusat'),
('Christine Kelly', 'christine.k@example.com', 'HR', 'HR Business Partner', '081234567955', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),
('Brandon Sanders', 'brandon.s@example.com', 'Developer', 'Backend Developer', '081234567956', 'Jl. H.R. Rasuna Said Kav. X-25, Jakarta Selatan'),
('Lauren Price', 'lauren.p@example.com', 'Manager', 'Product Manager', '081234567957', 'Jl. Jendral Sudirman No. 35, Jakarta Pusat'),
('Ethan Bennett', 'ethan.b@example.com', 'Developer', 'Frontend Developer', '081234567958', 'Jl. K.H. Wahid Hasyim No. 32, Jakarta Pusat'),
('Megan Wood', 'megan.w@example.com', 'QA', 'Test Automation Engineer', '081234567959', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),

-- Continue with more sample data...
('Christian Barnes', 'christian.b@example.com', 'Designer', 'UI Designer', '081234567960', 'Jl. H.R. Rasuna Said Kav. C-55, Jakarta Selatan'),
('Victoria Ross', 'victoria.r@example.com', 'HR', 'Recruiter', '081234567961', 'Jl. Jendral Sudirman No. 40, Jakarta Pusat'),
('Dylan Henderson', 'dylan.h@example.com', 'Developer', 'Full Stack Developer', '081234567962', 'Jl. K.H. Mas Mansyur No. 45, Jakarta Pusat'),
('Julia Coleman', 'julia.c@example.com', 'Manager', 'Delivery Manager', '081234567963', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),
('Cameron Jenkins', 'cameron.j@example.com', 'Developer', 'Backend Developer', '081234567964', 'Jl. H.R. Rasuna Said Kav. X-30, Jakarta Selatan'),
('Lauren Perry', 'lauren.p2@example.com', 'QA', 'QA Lead', '081234567965', 'Jl. Jendral Sudirman No. 45, Jakarta Pusat'),
('Nathan Powell', 'nathan.p@example.com', 'Designer', 'Product Designer', '081234567966', 'Jl. K.H. Wahid Hasyim No. 38, Jakarta Pusat'),
('Samantha Long', 'samantha.l@example.com', 'HR', 'HR Generalist', '081234567967', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Tyler Foster', 'tyler.f@example.com', 'Developer', 'Frontend Developer', '081234567968', 'Jl. H.R. Rasuna Said Kav. C-65, Jakarta Selatan'),
('Rachel Gonzales', 'rachel.g@example.com', 'Manager', 'Project Manager', '081234567969', 'Jl. Jendral Sudirman No. 50, Jakarta Pusat'),

-- Continue with more sample data...
('Kyle Bryant', 'kyle.b@example.com', 'Developer', 'Backend Developer', '081234567970', 'Jl. K.H. Mas Mansyur No. 50, Jakarta Pusat'),
('Hannah Russell', 'hannah.r@example.com', 'QA', 'Test Engineer', '081234567971', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Brandon Griffin', 'brandon.g@example.com', 'Designer', 'UI/UX Designer', '081234567972', 'Jl. H.R. Rasuna Said Kav. X-35, Jakarta Selatan'),
('Kayla Diaz', 'kayla.d@example.com', 'HR', 'HR Manager', '081234567973', 'Jl. Jendral Sudirman No. 55, Jakarta Pusat'),
('Aaron Hayes', 'aaron.h2@example.com', 'Developer', 'Full Stack Developer', '081234567974', 'Jl. K.H. Wahid Hasyim No. 42, Jakarta Pusat'),
('Natalie Myers', 'natalie.m@example.com', 'Manager', 'Technical Lead', '081234567975', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Justin Ford', 'justin.f@example.com', 'Developer', 'Mobile Developer', '081234567976', 'Jl. H.R. Rasuna Said Kav. C-75, Jakarta Selatan'),
('Christina Hamilton', 'christina.h@example.com', 'QA', 'QA Engineer', '081234567977', 'Jl. Jendral Sudirman No. 60, Jakarta Pusat'),
('Austin Graham', 'austin.g@example.com', 'Designer', 'UX Researcher', '081234567978', 'Jl. K.H. Mas Mansyur No. 55, Jakarta Pusat'),
('Allison Sullivan', 'allison.s@example.com', 'HR', 'HR Business Partner', '081234567979', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),

-- Continue with more sample data...
('Gabriel Wallace', 'gabriel.w@example.com', 'Developer', 'Backend Developer', '081234567980', 'Jl. H.R. Rasuna Said Kav. X-40, Jakarta Selatan'),
('Mariah Woods', 'mariah.w@example.com', 'Manager', 'Product Manager', '081234567981', 'Jl. Jendral Sudirman No. 65, Jakarta Pusat'),
('Cody Cole', 'cody.c@example.com', 'Developer', 'Frontend Developer', '081234567982', 'Jl. K.H. Wahid Hasyim No. 48, Jakarta Pusat'),
('Vanessa West', 'vanessa.w@example.com', 'QA', 'Test Automation Engineer', '081234567983', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Connor Jordan', 'connor.j@example.com', 'Designer', 'UI Designer', '081234567984', 'Jl. H.R. Rasuna Said Kav. C-85, Jakarta Selatan'),
('Alexis Owens', 'alexis.o@example.com', 'HR', 'Recruiter', '081234567985', 'Jl. Jendral Sudirman No. 70, Jakarta Pusat'),
('Caleb Reynolds', 'caleb.r@example.com', 'Developer', 'Full Stack Developer', '081234567986', 'Jl. K.H. Mas Mansyur No. 60, Jakarta Pusat'),
('Mackenzie Fisher', 'mackenzie.f@example.com', 'Manager', 'Delivery Manager', '081234567987', 'Jl. Jendral Sudirman Kav. 45, Jakarta Selatan'),
('Noah Ellis', 'noah.e@example.com', 'Developer', 'Backend Developer', '081234567988', 'Jl. H.R. Rasuna Said Kav. X-45, Jakarta Selatan'),
('Kaitlyn Harrison', 'kaitlyn.h@example.com', 'QA', 'QA Lead', '081234567989', 'Jl. Jendral Sudirman No. 75, Jakarta Pusat'),

-- Continue with more sample data...
('Lucas Gibson', 'lucas.g@example.com', 'Designer', 'Product Designer', '081234567990', 'Jl. K.H. Wahid Hasyim No. 52, Jakarta Pusat'),
('Haley Mcdonald', 'haley.m@example.com', 'HR', 'HR Generalist', '081234567991', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Evan Cruz', 'evan.c@example.com', 'Developer', 'Frontend Developer', '081234567992', 'Jl. H.R. Rasuna Said Kav. C-95, Jakarta Selatan'),
('Alexandra Marshall', 'alexandra.m@example.com', 'Manager', 'Project Manager', '081234567993', 'Jl. Jendral Sudirman No. 80, Jakarta Pusat'),
('Logan Ortiz', 'logan.o@example.com', 'Developer', 'Backend Developer', '081234567994', 'Jl. K.H. Mas Mansyur No. 65, Jakarta Pusat'),
('Brianna Gomez', 'brianna.g@example.com', 'QA', 'Test Engineer', '081234567995', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat'),
('Jared Murray', 'jared.m@example.com', 'Designer', 'UI/UX Designer', '081234567996', 'Jl. H.R. Rasuna Said Kav. X-50, Jakarta Selatan'),
('Jasmine Freeman', 'jasmine.f@example.com', 'HR', 'HR Manager', '081234567997', 'Jl. Jendral Sudirman No. 85, Jakarta Pusat'),
('Owen Wells', 'owen.w@example.com', 'Developer', 'Full Stack Developer', '081234567998', 'Jl. K.H. Wahid Hasyim No. 58, Jakarta Pusat'),
('Ariana Webb', 'ariana.w@example.com', 'Manager', 'Technical Lead', '081234567999', 'Jl. Jendral Sudirman Kav. 45, Jakarta Pusat');
