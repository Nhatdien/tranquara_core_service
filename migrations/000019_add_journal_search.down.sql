-- Rollback Migration 000019: Remove full-text search from user_journals

-- Drop the GIN index
DROP INDEX IF EXISTS idx_user_journals_search;

-- Drop the search_vector column
ALTER TABLE user_journals DROP COLUMN IF EXISTS search_vector;
