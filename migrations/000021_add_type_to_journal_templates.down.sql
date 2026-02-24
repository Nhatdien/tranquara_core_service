-- Rollback migration 000021: Remove type column from journal_templates

DROP INDEX IF EXISTS idx_journal_templates_type_category;
DROP INDEX IF EXISTS idx_journal_templates_type;
ALTER TABLE journal_templates DROP COLUMN IF EXISTS type;
