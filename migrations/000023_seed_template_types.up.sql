-- Migration 000023: Set type for existing journal_templates
-- Learning collections (educational doc-based content) → type = 'learn'
-- Journal collections (writing prompts) → type = 'journal' (default, no update needed)

-- Collection 2: Therapy Preparation (Learning - doc slides about therapy signs)
UPDATE journal_templates SET type = 'learn' WHERE id = '33333333-3333-3333-3333-333333333333';

-- Collection 4: Introduction to Journaling (Learning)
UPDATE journal_templates SET type = 'learn' WHERE id = '22222222-2222-2222-2222-222222222222';

-- Collection 5: Understanding Anxiety (Learning)
UPDATE journal_templates SET type = 'learn' WHERE id = 'bbbb2222-bbbb-4222-bbbb-bbbbbbbb2222';

-- Collection 6: Better Sleep (Learning)
UPDATE journal_templates SET type = 'learn' WHERE id = 'cccc3333-cccc-4333-cccc-cccccccc3333';

-- Collection 9: Understanding Emotions (Learning)
UPDATE journal_templates SET type = 'learn' WHERE id = 'ffff6666-ffff-4666-ffff-ffffffff6666';

-- Collection 10: Mindfulness (Learning)
UPDATE journal_templates SET type = 'learn' WHERE id = 'a0a07777-a0a0-4777-a0a0-a0a0a0a07777';

-- Collection 11: Self-Compassion (Learning)
UPDATE journal_templates SET type = 'learn' WHERE id = 'b0b08888-b0b0-4888-b0b0-b0b0b0b08888';

-- Remaining collections keep default type = 'journal':
-- Collection 1: Daily Reflection (55555555-...)
-- Collection 3: Stress Management (44444444-...)
-- Collection 7: Relationships & Connection (dddd4444-...)
-- Collection 8: Gratitude Practice (eeee5555-...)
-- Collection 12: Check-Ins (66666666-...)
