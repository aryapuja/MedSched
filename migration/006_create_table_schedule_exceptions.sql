CREATE TABLE schedule_exceptions (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER NOT NULL,
    exception_date DATE NOT NULL,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);