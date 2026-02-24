-- Migration 000022: Create user_learned_slide_groups table
-- Tracks which slide groups within learn-type collections a user has completed

CREATE TABLE user_learned_slide_groups (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    collection_id UUID NOT NULL REFERENCES journal_templates(id) ON DELETE CASCADE,
    slide_group_id VARCHAR(100) NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (user_id, collection_id, slide_group_id)
);

-- Indexes for common queries
CREATE INDEX idx_learned_user_id ON user_learned_slide_groups(user_id);
CREATE INDEX idx_learned_collection_id ON user_learned_slide_groups(collection_id);
CREATE INDEX idx_learned_user_collection ON user_learned_slide_groups(user_id, collection_id);

COMMENT ON TABLE user_learned_slide_groups IS 'Tracks per-slide-group completion for learn-type collections';
COMMENT ON COLUMN user_learned_slide_groups.slide_group_id IS 'Matches the id field inside journal_templates.slide_groups JSONB array';
