CREATE TABLE user_journals (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    title VARCHAR(50),
    short_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);