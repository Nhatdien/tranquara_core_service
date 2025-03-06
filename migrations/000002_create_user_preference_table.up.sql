CREATE TABLE user_preferences (
    user_id INT PRIMARY KEY REFERENCES users(user_id),
    program_mode VARCHAR(50) CHECK (program_mode IN ('8-week', 'self-guided')),
    daily_reminder_time TIME,
    notification_enabled BOOLEAN DEFAULT TRUE
);
