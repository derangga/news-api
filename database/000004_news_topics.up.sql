begin;

CREATE TABLE news_topics (
    id SERIAL PRIMARY KEY,
    news_article_id INTEGER NOT NULL REFERENCES news_articles(id) ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(news_article_id, topic_id)
);

CREATE INDEX idx_news_topics_news_id ON news_topics(news_article_id);
CREATE INDEX idx_news_topics_topic_id ON news_topics(topic_id);

commit;