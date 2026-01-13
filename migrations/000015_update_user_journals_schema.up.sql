-- Update user_journals table for offline-first journaling with embedded emotions/AI
-- Migration 000015: Add new fields for TipTap content and mood tracking

-- Drop old table since we're starting fresh (dev environment)
DROP TABLE IF EXISTS user_journals CASCADE;

-- Recreate with new schema
CREATE TABLE user_journals (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    collection_id UUID,                          -- Renamed from template_id
    title VARCHAR(255),
    content TEXT NOT NULL,                       -- TipTap JSON with embedded emotions + AI
    content_html TEXT,                           -- Rendered HTML preview for search/display
    mood_score INTEGER CHECK (mood_score >= 1 AND mood_score <= 10), -- Numeric 1-10
    mood_label VARCHAR(50),                      -- Weather label: "Storm", "Sunny", etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user_informations(user_id) ON DELETE CASCADE,
    FOREIGN KEY(collection_id) REFERENCES journal_templates(id) ON DELETE SET NULL
);

-- Create indexes for performance
CREATE INDEX idx_user_journals_user_id ON user_journals(user_id);
CREATE INDEX idx_user_journals_created_at ON user_journals(created_at DESC);
CREATE INDEX idx_user_journals_collection_id ON user_journals(collection_id);
CREATE INDEX idx_user_journals_mood_score ON user_journals(mood_score);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_journals_updated_at BEFORE UPDATE
    ON user_journals FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
