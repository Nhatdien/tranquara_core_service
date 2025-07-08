CREATE TABLE user_streaks (
    user_id UUID,
    current_streak INT DEFAULT 1,  
    longest_streak INT DEFAULT 1,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
);
