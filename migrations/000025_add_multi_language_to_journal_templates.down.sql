-- Rollback migration 000025: Remove multi-language columns from journal_templates

ALTER TABLE journal_templates
    DROP COLUMN IF EXISTS title_vi,
    DROP COLUMN IF EXISTS description_vi,
    DROP COLUMN IF EXISTS slide_groups_vi;
