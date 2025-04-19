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
INSERT INTO schedules (doctor_id, day_of_week, start_time, end_time, available_slot) 
VALUES 
(1, 'Senin', '09:00:00', '12:00:00', 2),
(1, 'Selasa', '13:00:00', '16:00:00', 2),
(2, 'Senin', '10:00:00', '14:00:00', 2),
(2, 'Rabu', '08:00:00', '12:00:00', 2);

INSERT INTO schedule_exceptions (schedule_id, exception_date) 
VALUES 
(1, '2025-04-21'),
(2, '2025-04-23');

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