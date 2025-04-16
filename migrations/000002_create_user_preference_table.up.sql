CREATE TABLE user_information (
    user_id INT PRIMARY KEY,
    age     SMALLINT,
    kyc_answers JSONB,
    program_mode VARCHAR(50) CHECK (program_mode IN ('8-week', 'self-guided')),
    daily_reminder_time TIME,
    notification_enabled BOOLEAN DEFAULT TRUE
);
