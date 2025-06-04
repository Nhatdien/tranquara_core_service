CREATE TABLE user_completed_exercises (
    id UUID DEFAULT gen_random_uuid() ,
    user_id UUID,
    exercise_id UUID REFERENCES exercises(exercise_id),
    duration SMALLINT CHECK(duration > 0),
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);