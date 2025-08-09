CREATE TABLE
    user_information (
        user_id UUID PRIMARY KEY,
        kyc_answers JSONB,
        name TEXT,
        age_range VARCHAR(50),
        gender VARCHAR(50),
        settings JSONB,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );