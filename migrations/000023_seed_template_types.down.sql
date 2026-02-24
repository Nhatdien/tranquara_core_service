-- Rollback: Reset all templates back to default 'journal' type
UPDATE journal_templates SET type = 'journal' WHERE type = 'learn';
