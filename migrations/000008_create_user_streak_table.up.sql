CREATE TABLE user_streaks (
    user_id INT,
    current_streak INT DEFAULT 0,  
    longest_streak INT DEFAULT 0,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
);
