CREATE TABLE user_information (
    user_id UUID PRIMARY KEY,
    kyc_answers JSONB,
    name    TEXT,
    age     SMALLINT,
    gender  VARCHAR(50),
    settings JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
