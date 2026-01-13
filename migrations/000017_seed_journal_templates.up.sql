-- Seed journal_templates (Collections) with sample data for testing
-- Migration 000017: Insert sample collections

-- Collection 1: Daily Reflection
INSERT INTO journal_templates (id, title, description, category, slide_groups, is_active) VALUES (
    '55555555-5555-5555-5555-555555555555'::uuid,
    'Daily Reflection',
    'Simple daily prompts to help you reflect on your mornings and evenings with intention.',
    'self_care',
    $$[
        {
            "id": "morning-prep",
            "title": "Morning",
            "description": "Start your day with mindful journaling and positive focus.",
            "position": 1,
            "slides": [
                {
                    "id": "morning-mood",
                    "type": "emotion_log",
                    "question": "How are you feeling this morning?",
                    "config": {
                        "scale": "1-10",
                        "labels": ["Storm", "Heavy Rain", "Rain", "Cloudy", "Partly Cloudy", "Mostly Sunny", "Sunny", "Bright", "Radiant", "Blissful"]
                    }
                },
                {
                    "id": "morning-sleep",
                    "type": "sleep_check",
                    "question": "How many hours did you sleep last night?",
                    "config": {
                        "min": 0,
                        "max": 12
                    }
                },
                {
                    "id": "morning-mind",
                    "type": "journal_prompt",
                    "question": "What is on my mind this morning?",
                    "config": {
                        "allowAI": true,
                        "minLength": 20
                    }
                },
                {
                    "id": "morning-intentions",
                    "type": "journal_prompt",
                    "question": "What can I do to make today amazing?",
                    "config": {
                        "allowAI": true,
                        "minLength": 20
                    }
                }
            ]
        },
        {
            "id": "evening-reflection",
            "title": "Evening",
            "description": "Close your day with reflection and self-kindness.",
            "position": 2,
            "slides": [
                {
                    "id": "evening-mood",
                    "type": "emotion_log",
                    "question": "How are you feeling this evening?",
                    "config": {
                        "scale": "1-10",
                        "labels": ["Storm", "Heavy Rain", "Rain", "Cloudy", "Partly Cloudy", "Mostly Sunny", "Sunny", "Bright", "Radiant", "Blissful"]
                    }
                },
                {
                    "id": "evening-memorable",
                    "type": "journal_prompt",
                    "question": "What happened today worth remembering?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "evening-positive",
                    "type": "journal_prompt",
                    "question": "What went well today?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "evening-learning",
                    "type": "journal_prompt",
                    "question": "What did I learn today?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "evening-changes",
                    "type": "journal_prompt",
                    "question": "What would I have changed about today?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "evening-celebrate",
                    "type": "journal_prompt",
                    "question": "What can I celebrate today?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "evening-growth",
                    "type": "journal_prompt",
                    "question": "How am I different from yesterday?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "evening-winddown",
                    "type": "journal_prompt",
                    "question": "How can I wind down, release the day, and rest now?",
                    "config": {
                        "allowAI": true
                    }
                }
            ]
        }
    ]$$::jsonb,
    true
);

-- Collection 2: Therapy Preparation
INSERT INTO journal_templates (id, title, description, category, slide_groups, is_active) VALUES (
    '33333333-3333-3333-3333-333333333333'::uuid,
    'Therapy Preparation',
    'Learn how to recognize when you may need therapy and what to expect in your first session.',
    'mental_health',
    $$[
        {
            "id": "signs-need-therapy",
            "title": "Signs that you need therapy",
            "description": "Learn to recognize when it may be time to seek professional support.",
            "position": 1,
            "slides": [
                {
                    "id": "sign-tearful",
                    "type": "doc",
                    "title": "You are tearful for no reason",
                    "content": "<h3>You are tearful for no reason</h3><p>When doing everyday activities like reading or watching a movie, if you find yourself tearful for no clear reason, it may be a sign your emotions are full and you need someone to talk to.</p>"
                },
                {
                    "id": "sign-negative",
                    "type": "doc",
                    "title": "You are always having negative emotions",
                    "content": "<h3>You are always having negative emotions</h3><p>Persistent negative emotions can weigh heavily and may signal that professional support could help.</p>"
                },
                {
                    "id": "sign-habits",
                    "type": "doc",
                    "title": "You slip back to unhealthy habits",
                    "content": "<h3>You slip back to unhealthy habits</h3><p>Falling into old unhealthy habits can be a sign of deeper struggles that therapy may help address.</p>"
                },
                {
                    "id": "sign-control",
                    "type": "doc",
                    "title": "You feel like your emotions are controlling you",
                    "content": "<h3>You feel like your emotions are controlling you</h3><p>Experiences like rage outbursts, yelling, or nonstop thoughts about past failures can interfere with sleep and daily life. This may be a signal that support is needed.</p>"
                },
                {
                    "id": "sign-disconnected",
                    "type": "doc",
                    "title": "You feel disconnected between your inner and outer world",
                    "content": "<h3>You feel disconnected between your inner and outer world</h3><p>You might seem cheerful in social settings but feel very different inside. Having separate faces socially and privately can indicate disconnection.</p>"
                },
                {
                    "id": "sign-relationships",
                    "type": "doc",
                    "title": "Your relationships are suffering",
                    "content": "<h3>Your relationships are suffering</h3><p>When you feel distant from friends or family, it can point to issues with trust, self-esteem, or deeper emotional struggles.</p>"
                },
                {
                    "id": "sign-reminder",
                    "type": "doc",
                    "title": "A reminder",
                    "content": "<h3>A reminder</h3><p>These are just common signs — you don't need to have any of them to see a therapist. All you need is the desire to feel better and share what is on your chest.</p>"
                },
                {
                    "id": "therapy-reflection",
                    "type": "journal_prompt",
                    "question": "After reading these signs, what resonates with you? What are you feeling right now?",
                    "config": {
                        "allowAI": true,
                        "minLength": 30
                    }
                }
            ]
        },
        {
            "id": "what-to-know",
            "title": "What I want my therapist to know",
            "description": "Reflect on what you would like to share in your first therapy session.",
            "position": 2,
            "slides": [
                {
                    "id": "therapist-know-intro",
                    "type": "doc",
                    "title": "Preparing for your first session",
                    "content": "<h3>Preparing for your first session</h3><p>Therapy can feel overwhelming at first. Taking time to reflect on what you want to share can help you feel more prepared and confident.</p><p>Use these prompts to explore what matters most to you right now.</p>"
                },
                {
                    "id": "therapist-main-concern",
                    "type": "journal_prompt",
                    "question": "What is the main thing bringing you to therapy right now?",
                    "config": {
                        "allowAI": true,
                        "minLength": 30
                    }
                },
                {
                    "id": "therapist-goals",
                    "type": "journal_prompt",
                    "question": "What do you hope to achieve through therapy?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "therapist-struggles",
                    "type": "journal_prompt",
                    "question": "What has been most difficult for you lately?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "therapist-support",
                    "type": "journal_prompt",
                    "question": "What kind of support do you think would be most helpful?",
                    "config": {
                        "allowAI": true
                    }
                }
            ]
        }
    ]$$::jsonb,
    true
);

-- Collection 3: Stress Management
INSERT INTO journal_templates (id, title, description, category, slide_groups, is_active) VALUES (
    '44444444-4444-4444-4444-444444444444'::uuid,
    'Stress Management',
    'Identify your stress sources and develop healthy coping strategies.',
    'self_care',
    $$[
        {
            "id": "stress-check-in",
            "title": "Stress Check-In",
            "description": "Assess your current stress level and identify triggers.",
            "position": 1,
            "slides": [
                {
                    "id": "stress-level",
                    "type": "emotion_log",
                    "question": "How stressed are you feeling right now?",
                    "config": {
                        "scale": "1-10",
                        "labels": ["Calm", "Relaxed", "Comfortable", "Slightly Tense", "Moderate", "Stressed", "Very Stressed", "Overwhelmed", "Extreme", "Breaking Point"]
                    }
                },
                {
                    "id": "stress-sources",
                    "type": "journal_prompt",
                    "question": "What is causing you stress right now?",
                    "config": {
                        "allowAI": true,
                        "minLength": 20
                    }
                },
                {
                    "id": "stress-body",
                    "type": "journal_prompt",
                    "question": "How is stress showing up in your body? (headaches, tension, fatigue, etc.)",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "stress-coping",
                    "type": "journal_prompt",
                    "question": "What has helped you cope with stress in the past?",
                    "config": {
                        "allowAI": true
                    }
                },
                {
                    "id": "stress-action",
                    "type": "journal_prompt",
                    "question": "What is one small thing you can do today to reduce your stress?",
                    "config": {
                        "allowAI": true
                    }
                }
            ]
        }
    ]$$::jsonb,
    true
);
