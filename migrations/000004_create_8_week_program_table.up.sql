CREATE TABLE program_exercises (
    id SERIAL PRIMARY KEY,
    week_number INT CHECK (week_number BETWEEN 1 AND 8),
    day_number INT CHECK (day_number BETWEEN 1 AND 7),
    exercise_id INT REFERENCES exercises(exercise_id),
    UNIQUE (week_number, day_number, exercise_id)
);