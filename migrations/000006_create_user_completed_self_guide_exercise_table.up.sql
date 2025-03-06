CREATE TABLE user_self_guided_activities (
    activity_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    exercise_id INT REFERENCES exercises(exercise_id) ON DELETE SET NULL,
    custom_activity_name VARCHAR(255), -- Allows user-created activities
    duration_minutes INT CHECK (duration_minutes >= 1),
    notes TEXT,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (exercise_id IS NOT NULL AND custom_activity_name IS NULL) OR
        (exercise_id IS NULL AND custom_activity_name IS NOT NULL)
    )
);
