CREATE TABLE user_informations (
    user_id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100),
    oauth_provider VARCHAR(50) DEFAULT 'email',
    kyc_answers JSONB,
    name TEXT,
    age_range VARCHAR(50),
    gender VARCHAR(50),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_informations_email ON user_informations(email);
CREATE INDEX idx_user_informations_oauth_provider ON user_informations(oauth_provider);
