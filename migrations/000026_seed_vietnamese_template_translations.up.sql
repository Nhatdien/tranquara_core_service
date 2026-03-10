-- Migration 000026: Add Vietnamese translations for template titles and descriptions
-- Sets title_vi and description_vi for all existing journal_templates

-- Collection 1: Daily Reflection (journal)
UPDATE journal_templates SET
    title_vi = 'Suy ngẫm hàng ngày',
    description_vi = 'Câu hỏi đơn giản hàng ngày giúp bạn suy ngẫm về buổi sáng và buổi tối một cách có ý thức.'
WHERE id = '55555555-5555-5555-5555-555555555555';

-- Collection 2: Therapy Preparation (learn)
UPDATE journal_templates SET
    title_vi = 'Chuẩn bị trị liệu',
    description_vi = 'Tìm hiểu cách nhận biết khi nào bạn cần trị liệu và những gì cần chuẩn bị cho buổi đầu tiên.'
WHERE id = '33333333-3333-3333-3333-333333333333';

-- Collection 3: Stress Management (journal)
UPDATE journal_templates SET
    title_vi = 'Quản lý căng thẳng',
    description_vi = 'Nhận diện nguồn gốc căng thẳng và phát triển chiến lược đối phó lành mạnh.'
WHERE id = '44444444-4444-4444-4444-444444444444';

-- Collection 4: Introduction to Journaling (learn)
UPDATE journal_templates SET
    title_vi = 'Giới thiệu về viết nhật ký',
    description_vi = 'Tìm hiểu những điều cơ bản về viết nhật ký, tại sao nó quan trọng và lợi ích cho sức khỏe tinh thần.'
WHERE id = '22222222-2222-2222-2222-222222222222';

-- Collection 5: Understanding Anxiety (learn)
UPDATE journal_templates SET
    title_vi = 'Hiểu về lo âu',
    description_vi = 'Tìm hiểu về lo âu, các yếu tố kích hoạt và kỹ thuật dựa trên bằng chứng để quản lý lo lắng và sợ hãi.'
WHERE id = 'bbbb2222-bbbb-4222-bbbb-bbbbbbbb2222';

-- Collection 6: Better Sleep (learn)
UPDATE journal_templates SET
    title_vi = 'Giấc ngủ tốt hơn',
    description_vi = 'Tìm hiểu về vệ sinh giấc ngủ và phát triển thói quen hỗ trợ giấc ngủ sâu, phục hồi sức khỏe.'
WHERE id = 'cccc3333-cccc-4333-cccc-cccccccc3333';

-- Collection 7: Relationships & Connection (journal)
UPDATE journal_templates SET
    title_vi = 'Mối quan hệ & Kết nối',
    description_vi = 'Khám phá các mối quan hệ, cách giao tiếp và cách xây dựng kết nối sâu sắc hơn với người khác.'
WHERE id = 'dddd4444-dddd-4444-dddd-dddddddd4444';

-- Collection 8: Gratitude Practice (journal)
UPDATE journal_templates SET
    title_vi = 'Thực hành biết ơn',
    description_vi = 'Nuôi dưỡng lòng biết ơn qua các bài suy ngẫm hàng ngày, chuyển sự chú ý sang sự trân trọng và tích cực.'
WHERE id = 'eeee5555-eeee-4555-eeee-eeeeeeee5555';

-- Collection 9: Understanding Emotions (learn)
UPDATE journal_templates SET
    title_vi = 'Hiểu về cảm xúc',
    description_vi = 'Tìm hiểu về cảm xúc, tại sao chúng quan trọng và cách làm việc cùng chúng thay vì chống lại.'
WHERE id = 'ffff6666-ffff-4666-ffff-ffffffff6666';

-- Collection 10: Mindfulness (learn)
UPDATE journal_templates SET
    title_vi = 'Chánh niệm',
    description_vi = 'Học và thực hành chánh niệm để nuôi dưỡng nhận thức khoảnh khắc hiện tại và sự bình yên nội tâm.'
WHERE id = 'a0a07777-a0a0-4777-a0a0-a0a0a0a07777';

-- Collection 11: Self-Compassion (learn)
UPDATE journal_templates SET
    title_vi = 'Tự thương yêu bản thân',
    description_vi = 'Học cách đối xử với bản thân bằng sự tử tế mà bạn dành cho một người bạn tốt.'
WHERE id = 'b0b08888-b0b0-4888-b0b0-b0b0b0b08888';

-- Collection 12: Check-Ins (journal)
UPDATE journal_templates SET
    title_vi = 'Ghi nhận cảm xúc',
    description_vi = 'Câu hỏi hướng dẫn cho việc ghi nhận hàng ngày và cảm xúc, giúp tăng nhận thức bản thân và cân bằng cảm xúc.'
WHERE id = '66666666-6666-4666-6666-666666666666';
