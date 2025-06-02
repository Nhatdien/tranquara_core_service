CREATE TABLE journal_templates (
    id UUID PRIMARY KEY,
    title VARCHAR(50),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)