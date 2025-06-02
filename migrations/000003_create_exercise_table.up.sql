CREATE TABLE exercises (
    exercise_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    media_link TEXT, 
    exercise_type VARCHAR(100)
);