CREATE TABLE forum_articles (
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT                      REFERENCES users,
  title      VARCHAR(255)                NOT NULL,  
  body       TEXT                        NOT NULL,
  type       VARCHAR(8)                  NOT NULL DEFAULT 'markdown',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_forum_articles
  ON forum_articles (title);
CREATE INDEX idx_forum_type
  ON forum_articles (type);

CREATE TABLE forum_tags (
  id         BIGSERIAL PRIMARY KEY,
  name       VARCHAR(255)                NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_forum_tags_name
  ON forum_tags (name);

CREATE TABLE forum_articles_tags (
  article_id BIGINT REFERENCES forum_articles ON DELETE CASCADE,
  tag_id     BIGINT REFERENCES forum_tags ON DELETE CASCADE,
  PRIMARY KEY (article_id, tag_id)
);

CREATE TABLE forum_comments (
  id         BIGSERIAL PRIMARY KEY,
  article_id BIGINT                      REFERENCES forum_articles,
  user_id    BIGINT                      REFERENCES users,
  body       TEXT                        NOT NULL,
  type       VARCHAR(8)                  NOT NULL DEFAULT 'markdown',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_forum_comments_type
  ON forum_comments (type);
