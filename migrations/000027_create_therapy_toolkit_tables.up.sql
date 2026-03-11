-- Therapy sessions
CREATE TABLE IF NOT EXISTS therapy_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    session_date TIMESTAMPTZ,
    status TEXT DEFAULT 'scheduled',
    mood_before INTEGER CHECK (mood_before BETWEEN 1 AND 10),
    talking_points TEXT,
    session_priority TEXT,
    prep_pack_id UUID,
    mood_after INTEGER CHECK (mood_after BETWEEN 1 AND 10),
    key_takeaways TEXT,
    session_rating INTEGER CHECK (session_rating BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Prep packs (created now so FK can reference it; populated in Phase 3)
CREATE TABLE IF NOT EXISTS prep_packs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    date_range_start DATE NOT NULL,
    date_range_end DATE NOT NULL,
    content JSONB NOT NULL,
    journal_count INTEGER DEFAULT 0,
    personal_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Homework items
CREATE TABLE IF NOT EXISTS homework_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES therapy_sessions(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- FK: therapy_sessions.prep_pack_id → prep_packs.id
ALTER TABLE therapy_sessions
    ADD CONSTRAINT fk_sessions_prep_pack
    FOREIGN KEY (prep_pack_id) REFERENCES prep_packs(id);

-- Indexes
CREATE INDEX idx_therapy_sessions_user ON therapy_sessions(user_id);
CREATE INDEX idx_therapy_sessions_date ON therapy_sessions(session_date);
CREATE INDEX idx_prep_packs_user ON prep_packs(user_id);
CREATE INDEX idx_homework_items_session ON homework_items(session_id);
CREATE INDEX idx_homework_items_user ON homework_items(user_id);
