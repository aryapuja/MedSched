CREATE TABLE Schedules (
    id SERIAL PRIMARY KEY,
    doctor_id INT REFERENCES Doctors(id) ON DELETE CASCADE,
    schedule_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    available_slots INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);