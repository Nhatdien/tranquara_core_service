CREATE TABLE IF NOT EXISTS ai_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'preferences',
    source_journal_ids UUID[] DEFAULT '{}',
    confidence REAL NOT NULL DEFAULT 0.5,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ai_memories_user_id ON ai_memories(user_id);
CREATE INDEX idx_ai_memories_category ON ai_memories(user_id, category);
CREATE INDEX idx_ai_memories_created ON ai_memories(user_id, created_at DESC);
