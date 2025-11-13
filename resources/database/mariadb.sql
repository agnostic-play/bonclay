-- MariaDB 10.5+ / 10.6+ compatible
SET NAMES utf8mb4;
SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE';

-- =========================================================
-- Table: collections
-- =========================================================
CREATE TABLE `collections` (
                               `id` CHAR(36) NOT NULL,
                               `squad_id` CHAR(36) NOT NULL,
                               `name` VARCHAR(75) NOT NULL,
                               `baseURL` VARCHAR(255) DEFAULT NULL,
                               `slug` VARCHAR(125) NOT NULL,
                               `desc` VARCHAR(255) DEFAULT NULL,
                               `docs` VARCHAR(255) DEFAULT NULL,
                               `forward_proxy_url` VARCHAR(255) DEFAULT NULL,
                               `is_proxy_enable` TINYINT(1) NOT NULL DEFAULT 0,
                               `created_at` TIMESTAMP NULL DEFAULT NULL,
                               `updated_at` TIMESTAMP NULL DEFAULT NULL,
                               `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                               PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Table: custom_variables
-- =========================================================
CREATE TABLE `custom_variables` (
                                    `id` CHAR(36) NOT NULL,
                                    `collection_id` CHAR(36) DEFAULT NULL,
                                    `key` VARCHAR(255) NOT NULL,
                                    `value` TEXT DEFAULT NULL,
                                    `created_at` TIMESTAMP NULL DEFAULT NULL,
                                    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                    `updated_at` TIMESTAMP NULL DEFAULT NULL,
                                    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Table: diagrams
-- =========================================================
CREATE TABLE `diagrams` (
                            `id` CHAR(36) NOT NULL,
                            `collection_id` CHAR(36) NOT NULL,
                            `title` VARCHAR(255) NOT NULL,
                            `description` TEXT DEFAULT NULL,
                            `syntax_type` VARCHAR(50) DEFAULT 'mermaid',
                            `syntax` TEXT NOT NULL,
                            `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP(),
                            `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
                            `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Table: diagram_collections
-- =========================================================
CREATE TABLE `diagram_collections` (
                                       `id` CHAR(36) NOT NULL,
                                       `name` VARCHAR(255) NOT NULL,
                                       `description` TEXT DEFAULT NULL,
                                       `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP(),
                                       `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
                                       `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Table: endpoints
-- =========================================================
CREATE TABLE `endpoints` (
                             `id` CHAR(36) NOT NULL,
                             `category` VARCHAR(75) DEFAULT NULL,
                             `collection_id` CHAR(36) NOT NULL,
                             `path` VARCHAR(75) NOT NULL,
                             `desc` VARCHAR(120) DEFAULT NULL,
                             `method` ENUM('GET','POST','PUT','PATCH','DELETE','OPTIONS','TRACE') NOT NULL,
                             `active_scenario` CHAR(36) DEFAULT NULL,
                             `enable_response_intercept` TINYINT(1) NOT NULL DEFAULT 0,
                             `created_at` TIMESTAMP NULL DEFAULT NULL,
                             `updated_at` TIMESTAMP NULL DEFAULT NULL,
                             `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                             PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Table: scenario_response
-- =========================================================
CREATE TABLE `scenario_response` (
                                     `id` CHAR(36) NOT NULL,
                                     `endpoint_id` CHAR(36) NOT NULL,
                                     `desc` VARCHAR(120) DEFAULT NULL,
                                     `status_header` INT(11) DEFAULT NULL,
                                     `delay` TINYINT(4) DEFAULT NULL,
                                     `body` TEXT DEFAULT NULL,
                                     `header` LONGTEXT DEFAULT NULL,
                                     `created_at` TIMESTAMP NULL DEFAULT NULL,
                                     `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                     `updated_at` TIMESTAMP NULL DEFAULT NULL,
                                     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Table: squads
-- =========================================================
CREATE TABLE `squads` (
                          `id` CHAR(36) NOT NULL,
                          `name` VARCHAR(75) NOT NULL,
                          `slug` VARCHAR(125) NOT NULL,
                          `password` VARCHAR(255) NOT NULL,
                          `desc` VARCHAR(255) DEFAULT NULL,
                          `created_at` TIMESTAMP NULL DEFAULT NULL,
                          `updated_at` TIMESTAMP NULL DEFAULT NULL,
                          `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                          PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================
-- Triggers to auto-generate UUIDs
-- =========================================================
DELIMITER //

CREATE TRIGGER `bi_diagrams`
    BEFORE INSERT ON `diagrams`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER `bi_diagram_collections`
    BEFORE INSERT ON `diagram_collections`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER `bi_collections`
    BEFORE INSERT ON `collections`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER `bi_custom_variables`
    BEFORE INSERT ON `custom_variables`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER `bi_endpoints`
    BEFORE INSERT ON `endpoints`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER `bi_scenario_response`
    BEFORE INSERT ON `scenario_response`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER `bi_squads`
    BEFORE INSERT ON `squads`
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
    SET NEW.id = UUID();
END IF;
END//

DELIMITER ;
