-- +goose Up

CREATE TABLE IF NOT EXISTS feeds (
  id UUID PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  -- 确保同一个用户不能添加相同名称的订阅源
  UNIQUE (user_id, name),
  -- 确保同一个用户不能添加相同URL的订阅源
  UNIQUE (user_id, url)
);

-- +goose Down
DROP TABLE IF EXISTS feeds;