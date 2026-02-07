-- Rollback: Delete additional journal templates
DELETE FROM journal_templates WHERE id IN (
    '22222222-2222-2222-2222-222222222222',  -- Introduction to Journaling
    'bbbb2222-bbbb-4222-bbbb-bbbbbbbb2222',  -- Understanding Anxiety
    'cccc3333-cccc-4333-cccc-cccccccc3333',  -- Better Sleep
    'dddd4444-dddd-4444-dddd-dddddddd4444',  -- Relationships & Connection
    'eeee5555-eeee-4555-eeee-eeeeeeee5555',  -- Gratitude Practice
    'ffff6666-ffff-4666-ffff-ffffffff6666',  -- Understanding Emotions
    'a0a07777-a0a0-4777-a0a0-a0a0a0a07777',  -- Mindfulness
    'b0b08888-b0b0-4888-b0b0-b0b0b0b08888',  -- Self-Compassion
    '66666666-6666-4666-6666-666666666666'   -- Check-Ins
);
