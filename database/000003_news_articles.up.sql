begin;

CREATE TYPE article_status AS ENUM ('draft', 'published', 'deleted');

CREATE TABLE news_articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    status article_status NOT NULL DEFAULT 'draft',
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_news_articles_status ON news_articles(status);
CREATE INDEX idx_news_articles_published_at ON news_articles(published_at);
CREATE INDEX idx_news_articles_slug ON news_articles(slug);

commit;