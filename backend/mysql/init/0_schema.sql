CREATE TABLE IF NOT EXISTS posts (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  original_message TEXT NOT NULL,
  converted_message TEXT NOT NULL,
  user_name VARCHAR(32) NOT NULL,
  parent_id VARCHAR(36),
  root_id VARCHAR(36),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (parent_id) REFERENCES posts(id),
  FOREIGN KEY (root_id) REFERENCES posts(id)
);

CREATE TABLE IF NOT EXISTS posts_reactions (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  user_name VARCHAR(32) NOT NULL,
  reaction_id INT NOT NULL,
  post_id VARCHAR(36) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id),
  UNIQUE KEY (user_name, post_id, reaction_id)
);