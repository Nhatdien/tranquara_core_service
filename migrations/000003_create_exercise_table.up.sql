CREATE TABLE exercises (
    exercise_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INT CHECK (duration_minutes > 0),
    exercise_type VARCHAR(100)
);