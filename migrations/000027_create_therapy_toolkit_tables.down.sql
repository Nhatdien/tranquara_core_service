DROP INDEX IF EXISTS idx_homework_items_user;
DROP INDEX IF EXISTS idx_homework_items_session;
DROP INDEX IF EXISTS idx_prep_packs_user;
DROP INDEX IF EXISTS idx_therapy_sessions_date;
DROP INDEX IF EXISTS idx_therapy_sessions_user;

DROP TABLE IF EXISTS homework_items;

ALTER TABLE therapy_sessions DROP CONSTRAINT IF EXISTS fk_sessions_prep_pack;

DROP TABLE IF EXISTS therapy_sessions;
DROP TABLE IF EXISTS prep_packs;
