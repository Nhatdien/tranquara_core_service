-- Migration 000025: Add multi-language support to journal_templates
-- Adds Vietnamese (vi) columns for title, description, and slide_groups

ALTER TABLE journal_templates
    ADD COLUMN IF NOT EXISTS title_vi VARCHAR(255),
    ADD COLUMN IF NOT EXISTS description_vi TEXT,
    ADD COLUMN IF NOT EXISTS slide_groups_vi JSONB;

-- Add comment for documentation
COMMENT ON COLUMN journal_templates.title_vi IS 'Vietnamese translation of the template title';
COMMENT ON COLUMN journal_templates.description_vi IS 'Vietnamese translation of the template description';
COMMENT ON COLUMN journal_templates.slide_groups_vi IS 'Vietnamese translation of slide groups JSONB (same structure as slide_groups but with translated content)';
