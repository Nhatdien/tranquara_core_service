-- Rollback: Remove Vietnamese translations from template titles and descriptions
UPDATE journal_templates SET title_vi = NULL, description_vi = NULL;
