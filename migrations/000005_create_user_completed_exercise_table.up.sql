CREATE TABLE user_completed_exercises (
    id SERIAL PRIMARY KEY,
    user_id UUID,
    week_number INT CHECK (week_number BETWEEN 1 AND 8),
    day_number INT CHECK (day_number BETWEEN 1 AND 7),
    exercise_id INT REFERENCES exercises(exercise_id),
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);