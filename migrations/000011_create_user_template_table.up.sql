CREATE TABLE user_templates (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    title VARCHAR(50),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);