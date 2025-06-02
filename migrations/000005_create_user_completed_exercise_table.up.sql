CREATE TABLE user_completed_exercises (
    id SERIAL PRIMARY KEY,
    user_id UUID,
    exercise_id INT REFERENCES exercises(exercise_id),
    duration SMALLINT CHECK(duration > 0),
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);