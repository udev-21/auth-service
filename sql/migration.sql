
CREATE TABLE IF NOT EXISTS`users` (
  `id` char(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  `password` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime,  
  `created_by` char(36) DEFAULT NULL,
  `position` bigint AUTO_INCREMENT,
  `additional` json,
   KEY `position` (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE if not exists `service_owners`  (
  `id` char(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
  `user_id` char(36) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime,
  `created_by` char(36) DEFAULT NULL,
  `position` bigint AUTO_INCREMENT,
   KEY `position` (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- TODO: add foreign key to owner_id field 
CREATE TABLE if not exists `services`  (
  `id` char(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
  `owner_id` char(36) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(2048) DEFAULT NULL,
  `created_by` char(36) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime,  
  `position` bigint AUTO_INCREMENT,
   KEY `position` (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE if not exists `service_users`  (
  `id` char(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
  `user_id` char(36) NOT NULL,
  `service_id` char(36) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime,  
  `created_by` char(36) DEFAULT NULL,
  `position` bigint AUTO_INCREMENT,
   KEY `position` (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

