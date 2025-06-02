CREATE TABLE emotion_logs (
    id UUID PRIMARY KEY,
    user_id UUID,
    emotion VARCHAR(50),
    source VARCHAR(50),
    context TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)