CREATE TABLE IF NOT EXISTS alert_rules (
    id CHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    level VARCHAR(10) NOT NULL,
    disease VARCHAR(150) NOT NULL,
    description VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_system BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
    user_id CHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    role VARCHAR(20) DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS baby (
    baby_id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL UNIQUE,
    birth_date DATE NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_baby_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token VARCHAR(512) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    expires_at DATETIME(3) NOT NULL,
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sleep_sessions (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    sleep_time DATETIME(3) NOT NULL,
    wake_time DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    is_backdate BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_sleep_sessions_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_sleep_sessions_user_id ON sleep_sessions(user_id);
CREATE INDEX idx_sleep_sessions_sleep_time ON sleep_sessions(sleep_time);

CREATE TABLE IF NOT EXISTS symptom_logs (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    date DATETIME(3) NOT NULL,
    is_backdate BOOLEAN DEFAULT FALSE,
    physical_data TEXT NULL,
    last_alert_level VARCHAR(20) NULL,
    last_alert_data LONGTEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_symptom_logs_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_symptom_logs_user_id ON symptom_logs(user_id);
CREATE INDEX idx_symptom_logs_date ON symptom_logs(date);

CREATE TABLE IF NOT EXISTS bleeding_logs (
    id CHAR(36) PRIMARY KEY,
    log_id CHAR(36) NOT NULL,
    pad_usage VARCHAR(20) NULL,
    clot_size VARCHAR(20) NULL,
    blood_color VARCHAR(20) NULL,
    smell VARCHAR(20) NULL,
    CONSTRAINT fk_bleeding_logs_symptom_log FOREIGN KEY (log_id) REFERENCES symptom_logs(id) ON DELETE CASCADE
);

CREATE INDEX idx_bleeding_logs_log_id ON bleeding_logs(log_id);

CREATE TABLE IF NOT EXISTS mood_logs (
    id CHAR(36) PRIMARY KEY,
    log_id CHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    CONSTRAINT fk_mood_logs_symptom_log FOREIGN KEY (log_id) REFERENCES symptom_logs(id) ON DELETE CASCADE
);

CREATE INDEX idx_mood_logs_log_id ON mood_logs(log_id);

CREATE TABLE IF NOT EXISTS psychologists (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    title VARCHAR(100) NOT NULL,
    job VARCHAR(100) NOT NULL,
    experience_yr BIGINT NOT NULL,
    price_idr BIGINT NOT NULL,
    photo_url VARCHAR(255) NULL
);

CREATE TABLE IF NOT EXISTS psychologist_schedules (
    id CHAR(36) PRIMARY KEY,
    psychologist_id CHAR(36) NOT NULL,
    day_of_week VARCHAR(20) NOT NULL,
    start_time VARCHAR(10) NOT NULL,
    end_time VARCHAR(10) NOT NULL,
    CONSTRAINT fk_psychologist_schedules_psychologist FOREIGN KEY (psychologist_id) REFERENCES psychologists(id) ON DELETE CASCADE
);

CREATE INDEX idx_psychologist_schedules_psychologist_id ON psychologist_schedules(psychologist_id);
