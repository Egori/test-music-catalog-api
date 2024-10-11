CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,     
    text TEXT,
    link VARCHAR(255),
    release_date DATE,               
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Уникальный индекс для предотвращения дублирования песен
CREATE UNIQUE INDEX idx_group_title ON songs (group_name, title);

-- Индексы для ускорения запросов
CREATE INDEX idx_group_name ON songs (group_name);
CREATE INDEX idx_title ON songs (title);
CREATE INDEX idx_release_date ON songs (release_date);
