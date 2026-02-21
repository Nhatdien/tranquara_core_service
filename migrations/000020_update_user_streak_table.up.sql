-- Add missing columns from schema spec and fix defaults
ALTER TABLE user_streaks 
  ADD COLUMN IF NOT EXISTS total_entries INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Fix defaults: spec says DEFAULT 0, migration 000008 had DEFAULT 1
ALTER TABLE user_streaks ALTER COLUMN current_streak SET DEFAULT 0;
ALTER TABLE user_streaks ALTER COLUMN longest_streak SET DEFAULT 0;
