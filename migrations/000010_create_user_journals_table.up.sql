CREATE TABLE user_journals (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    template_id UUID,
    status VARCHAR(50) DEFAULT 'draft',
    title VARCHAR(50),
    short_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);