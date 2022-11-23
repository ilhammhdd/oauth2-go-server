CREATE TABLE IF NOT EXISTS `scopes`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `scope` TEXT NOT NULL,
  `permission` TINYINT UNSIGNED,
  `parent_scopes_id` BIGINT UNSIGNED,
  `clients_id` BIGINT UNSIGNED,
  FOREIGN KEY (`parent_scopes_id`) REFERENCES `scopes` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (`clients_id`) REFERENCES `clients` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `authz_codes`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `authz_code_id` TEXT NOT NULL UNIQUE KEY,
  `status` ENUM('authz_code_issued', 'authz_code_exchanged') NOT NULL,
  `status_at` DATETIME(6) NOT NULL,
  `clients_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`clients_id`) REFERENCES `clients` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `authz_codes_scopes`(
  `scopes_id` BIGINT UNSIGNED NOT NULL,
  `authz_codes_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`scopes_id`) REFERENCES `scopes` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (`authz_codes_id`) REFERENCES `authz_codes` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  PRIMARY KEY (`scopes_id`, `authz_codes_id`)
);

CREATE TABLE IF NOT EXISTS `access_xor_refresh_tokens`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `token_type` ENUM('access_token', 'refresh_token') NOT NULL,
  `token_id` TEXT NOT NULL,
  `status` ENUM('token_active', 'token_revoked') NOT NULL,
  `status_at` DATETIME(6) NOT NULL,
  `authz_codes_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`authz_codes_id`) REFERENCES `authz_codes` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);