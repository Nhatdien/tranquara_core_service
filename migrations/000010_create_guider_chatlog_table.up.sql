CREATE TABLE ai_guider_chatlog (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,  
    sender_type VARCHAR(50),
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
