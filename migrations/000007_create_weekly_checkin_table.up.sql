CREATE TABLE user_weekly_reflections (
    reflection_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    week_number INT CHECK (week_number BETWEEN 1 AND 8),
    user_notes TEXT,
    ai_insight TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);