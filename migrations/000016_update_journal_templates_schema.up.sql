-- Update journal_templates table to use JSONB slide_groups
-- Migration 000016: Convert templates to Collections with structured slides

-- Drop old table since we're starting fresh (dev environment)
DROP TABLE IF EXISTS journal_templates CASCADE;

-- Recreate with new Collections-based schema
CREATE TABLE journal_templates (
    id UUID DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),                    -- e.g., "Daily Reflection", "Therapy Prep"
    slide_groups JSONB NOT NULL,             -- Structured slide group definitions
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Create indexes
CREATE INDEX idx_journal_templates_category ON journal_templates(category);
CREATE INDEX idx_journal_templates_active ON journal_templates(is_active);
CREATE INDEX idx_journal_templates_slide_groups ON journal_templates USING GIN (slide_groups);

-- Create updated_at trigger for journal_templates
CREATE TRIGGER update_journal_templates_updated_at BEFORE UPDATE
    ON journal_templates FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comment for documentation
COMMENT ON COLUMN journal_templates.slide_groups IS 'JSONB array of slide group objects with id, title, description, and slides array';
