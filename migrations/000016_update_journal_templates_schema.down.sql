-- Rollback migration 000016
DROP TRIGGER IF EXISTS update_journal_templates_updated_at ON journal_templates;
DROP TABLE IF EXISTS journal_templates CASCADE;

-- Restore old schema
CREATE TABLE journal_templates (
    id UUID DEFAULT gen_random_uuid(),
    title VARCHAR(50),
    content VARCHAR(255)[],
    category VARCHAR(50),
    greetings TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
