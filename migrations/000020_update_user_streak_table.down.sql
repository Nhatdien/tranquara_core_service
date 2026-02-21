ALTER TABLE user_streaks DROP COLUMN IF EXISTS total_entries;
ALTER TABLE user_streaks DROP COLUMN IF EXISTS updated_at;

ALTER TABLE user_streaks ALTER COLUMN current_streak SET DEFAULT 1;
ALTER TABLE user_streaks ALTER COLUMN longest_streak SET DEFAULT 1;
