-- Migration 000021: Add type column to journal_templates
-- Distinguishes between 'journal' (journaling prompts) and 'learn' (micro-learning) collections

ALTER TABLE journal_templates ADD COLUMN type VARCHAR(50) NOT NULL DEFAULT 'journal';

CREATE INDEX idx_journal_templates_type ON journal_templates(type);

-- Composite index for filtering by type + category
CREATE INDEX idx_journal_templates_type_category ON journal_templates(type, category);

COMMENT ON COLUMN journal_templates.type IS 'Collection type: journal (prompts) or learn (educational)';
