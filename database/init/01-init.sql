-- データベース作成
CREATE DATABASE guideforge;

-- 使用するデータベースに接続
\c guideforge;

-- ユーザーテーブル
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  profile_image VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- マニュアルテーブル
CREATE TABLE manuals (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100),
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  is_public BOOLEAN DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 手順テーブル
CREATE TABLE steps (
  id SERIAL PRIMARY KEY,
  manual_id INTEGER NOT NULL REFERENCES manuals(id) ON DELETE CASCADE,
  order_number INTEGER NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 画像テーブル
CREATE TABLE images (
  id SERIAL PRIMARY KEY,
  step_id INTEGER REFERENCES steps(id) ON DELETE CASCADE,
  file_path VARCHAR(255) NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  file_size INTEGER NOT NULL,
  mime_type VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- インデックス作成
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_manuals_user_id ON manuals (user_id);
CREATE INDEX idx_manuals_category ON manuals (category);
CREATE INDEX idx_steps_manual_id ON steps (manual_id);
CREATE INDEX idx_steps_order_number ON steps (order_number);
CREATE INDEX idx_images_step_id ON images (step_id);

-- テスト用データの挿入（開発環境用）
INSERT INTO users (username, email, password_hash, created_at, updated_at)
VALUES ('testuser', 'test@example.com', '$2a$10$1qAz2wSx3eDc4rFv5tGb5edva6NKx.IfeanyP8w7VUZh5XILcOH.e', NOW(), NOW());

-- パスワードはbcryptハッシュ化された 'password' 