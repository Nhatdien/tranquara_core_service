-- Migration 000019: Add full-text search to user_journals
-- Uses PostgreSQL tsvector for efficient text search across title and content_html

-- Add generated tsvector column for full-text search
-- Weight 'A' for title (higher priority) and 'B' for content_html (lower priority)
ALTER TABLE user_journals 
ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(content_html, '')), 'B')
) STORED;

-- Create GIN index for fast full-text search
CREATE INDEX idx_user_journals_search ON user_journals USING GIN(search_vector);

-- Add comment for documentation
COMMENT ON COLUMN user_journals.search_vector IS 'Full-text search vector: title (weight A) + content_html (weight B)';
