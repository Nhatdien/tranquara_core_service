CREATE TABLE emotion_logs (
    id UUID DEFAULT gen_random_uuid() ,
    user_id UUID,
    emotion VARCHAR(50),
    source VARCHAR(50),
    context TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
)