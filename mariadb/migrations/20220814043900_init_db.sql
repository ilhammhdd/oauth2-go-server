CREATE TABLE IF NOT EXISTS `client_registrations`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `init_client_id_checksum` TEXT NOT NULL UNIQUE KEY,
  `basepoint` TEXT NOT NULL,
  `server_sk` TEXT NOT NULL,
  `server_pk` TEXT NOT NULL,
  `session_expired_at` DATETIME(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS `clients`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `token_endpoint_auth_method` SET('none', 'client_secret_post', 'client_secret_basic', 'client_secret_bearer'),
  `grant_types` SET('authorization_code', 'implicit', 'password', 'client_credentials', 'refresh_token'),
  `response_types` SET('code', 'token'),
  `client_name` TEXT NOT NULL,
  `client_uri` TEXT NOT NULL,
  `logo_uri` TEXT NOT NULL,
  `scope` TEXT NOT NULL,
  `tos_uri` TEXT NOT NULL,
  `policy_uri` TEXT NOT NULL,
  `software_id` TEXT NOT NULL,
  `software_version` TEXT NOT NULL,
  `init_client_id_checksum` TEXT NOT NULL UNIQUE KEY,
  `client_id` TEXT NOT NULL UNIQUE KEY,
  `client_id_issued_at` DATETIME(6) NOT NULL,
  `client_secret` TEXT NOT NULL,
  `client_secret_expired_at` DATETIME(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS `redirect_uris`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `uri` TEXT,
  `clients_id` BIGINT UNSIGNED NOT NULL,
  FOREIGN KEY (`clients_id`) REFERENCES `clients` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `contacts`(
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `soft_deleted_at` DATETIME(6),
  `contact` TEXT,
  `clients_id` BIGINT UNSIGNED,
  FOREIGN KEY (`clients_id`) REFERENCES `clients` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);