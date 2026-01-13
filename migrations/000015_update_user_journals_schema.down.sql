-- Rollback migration 000015
DROP TRIGGER IF EXISTS update_user_journals_updated_at ON user_journals;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS user_journals CASCADE;

-- Restore old schema
CREATE TABLE user_journals (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    template_id UUID,
    title VARCHAR(50),
    content TEXT, 
    mood VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);
