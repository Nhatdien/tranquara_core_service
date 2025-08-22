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