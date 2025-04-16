-- users
INSERT INTO Users (username, password, email, role, full_name, gender, birth_date) 
VALUES 
('alveria', 'hashed_password_123', 'alveria@example.com', 'pasien', 'Alveria', 'perempuan', '2002-02-04'),
('admin', 'hashed_password_456', 'admin@example.com', 'admin', 'admin', 'laki-laki', '1985-07-22');

-- doctor
INSERT INTO Doctors (name, specialization, contact_info)
VALUES 
('Dr. Tantri', 'Umum', '08123456789'),
('Dr. Sekar', 'Gigi', '08234567890');

-- schedule
INSERT INTO Schedules (doctor_id, schedule_date, start_time, end_time, available_slots)
VALUES 
(1, '2025-04-17', '09:00:00', '12:00:00', 5),
(1, '2025-04-17', '14:00:00', '17:00:00', 3),
(2, '2025-04-18', '08:00:00', '11:00:00', 4);

-- appointment
INSERT INTO Appointments (user_id, schedule_id, appointment_date, status)
VALUES 
(1, 1, '2025-04-17', 'pending'),
(2, 2, '2025-04-17', 'confirmed');

-- notification
INSERT INTO Notifications (user_id, appointment_id, notification_type, status)
VALUES 
(1, 1, 'reminder', 'sent'),
(2, 2, 'reminder', 'sent');