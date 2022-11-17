CREATE TABLE IF NOT EXISTS `url_one_time_tokens`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `pk` TEXT NOT NULL,
  `sk` TEXT NOT NULL,
  `one_time_token` TEXT NOT NULL,
  `signature` TEXT NOT NULL,
  `url` TEXT NOT NULL,
  `clients_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`clients_id`) REFERENCES `clients` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `email` TEXT NOT NULL,
  `password` TEXT NOT NULL,
  `clients_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`clients_id`) REFERENCES `clients` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `user_password_params`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `rand_salt` TEXT NOT NULL,
  `time` INT UNSIGNED NOT NULL,
  `memory` INT UNSIGNED NOT NULL,
  `threads` TINYINT(8) UNSIGNED NOT NULL,
  `keyLen` INT UNSIGNED NOT NULL,
  `users_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`users_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `usernames`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `username` TEXT NOT NULL,
  `unq_num` SMALLINT UNSIGNED NOT NULL,
  `users_id` BIGINT UNSIGNED NOT NULL,
  UNIQUE KEY `username_unq_num` (`username`, `unq_num`),
  FOREIGN KEY (`users_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS `username_idx` ON `usernames` (`username`);