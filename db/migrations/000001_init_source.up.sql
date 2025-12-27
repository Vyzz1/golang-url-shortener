CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    click_count BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT true
);


CREATE INDEX idx_short_code ON urls (short_code);
CREATE INDEX idx_created_at  ON urls (created_at);
CREATE INDEX idx_expires_at ON urls (expires_at);

CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    url_id BIGINT REFERENCES urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    referer TEXT
   
);

CREATE INDEX idx_ip_address ON clicks (ip_address);

CREATE INDEX idx_clicked_at ON clicks (clicked_at);
CREATE INDEX idx_url_id ON clicks (url_id);