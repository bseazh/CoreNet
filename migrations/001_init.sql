CREATE TABLE IF NOT EXISTS file (
  id         VARCHAR(32) PRIMARY KEY,
  owner_id   VARCHAR(32) NOT NULL,
  name       VARCHAR(255) NOT NULL,
  size       BIGINT NOT NULL,
  mime       VARCHAR(128),
  sha1       CHAR(40),
  version    INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS user_file (
  user_id   VARCHAR(32) NOT NULL,
  file_id   VARCHAR(32) NOT NULL,
  path      VARCHAR(1024) NOT NULL,
  tags      JSON NULL,
  acl       JSON NULL,
  PRIMARY KEY(user_id, file_id)
);

CREATE TABLE IF NOT EXISTS file_version (
  file_id     VARCHAR(32) NOT NULL,
  v           INT NOT NULL,
  object_uri  VARCHAR(1024) NOT NULL,
  PRIMARY KEY(file_id, v)
);

CREATE TABLE IF NOT EXISTS task (
  id         VARCHAR(32) PRIMARY KEY,
  type       VARCHAR(32) NOT NULL,   -- ocr|transcode
  file_id    VARCHAR(32) NOT NULL,
  status     VARCHAR(32) NOT NULL,   -- created|running|succeeded|failed
  result_uri VARCHAR(1024),
  err_msg    VARCHAR(1024),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
