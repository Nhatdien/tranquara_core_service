CREATE TABLE journal_templates (
    id UUID DEFAULT gen_random_uuid(),
    title VARCHAR(50),
    content TEXT,
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)