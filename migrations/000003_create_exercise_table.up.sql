CREATE TABLE exercises (
    exercise_id UUID DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    media_link TEXT, 
    exercise_type VARCHAR(100),

    PRIMARY KEY (exercise_id)
);