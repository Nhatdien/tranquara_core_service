CREATE TABLE user_journals (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    title VARCHAR(50),
    content TEXT,
    template_id UUID REFERENCES journal_templates(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);