CREATE TABLE user_self_guided_activities (
    activity_id SERIAL PRIMARY KEY,
    user_id UUID,
    exercise_id INT REFERENCES exercises(exercise_id) ON DELETE SET NULL,
    duration_minutes INT CHECK (duration_minutes >= 1),
    notes TEXT,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
