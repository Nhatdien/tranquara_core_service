CREATE TABLE user_journals (
    id UUID PRIMARY KEY,
    user_id UUID,
    title VARCHAR(50),
    content TEXT,
    template_id UUID REFERENCES journal_templates(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);