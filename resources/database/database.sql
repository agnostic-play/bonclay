-- =========================================
-- GLOBAL SETTINGS
-- =========================================
SET NAMES utf8mb4;
SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE';

-- =========================================
-- TABLE: squads
-- =========================================
CREATE TABLE `squads` (
                          `id` CHAR(36) NOT NULL,
                          `name` VARCHAR(75) NOT NULL,
                          `slug` VARCHAR(125) NOT NULL,
                          `password` VARCHAR(255) NOT NULL,
                          `desc` VARCHAR(255),
                          `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                          `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `uk_squads_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TABLE: collections
-- =========================================
CREATE TABLE `collections` (
                               `id` CHAR(36) NOT NULL,
                               `squad_id` CHAR(36) NOT NULL,
                               `name` VARCHAR(75) NOT NULL,
                               `base_url` VARCHAR(255),
                               `slug` VARCHAR(125) NOT NULL,
                               `desc` VARCHAR(255),
                               `docs` VARCHAR(255),
                               `forward_proxy_url` VARCHAR(255),
                               `is_proxy_enable` TINYINT(1) NOT NULL DEFAULT 0,
                               `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                               `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                               PRIMARY KEY (`id`),
                               KEY `idx_collections_squad` (`squad_id`),
                               UNIQUE KEY `uk_collections_slug` (`slug`),
                               CONSTRAINT `fk_collections_squad`
                                   FOREIGN KEY (`squad_id`) REFERENCES `squads`(`id`)
                                       ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TABLE: custom_variables
-- =========================================
CREATE TABLE `custom_variables` (
                                    `id` CHAR(36) NOT NULL,
                                    `collection_id` CHAR(36),
                                    `key` VARCHAR(255) NOT NULL,
                                    `value` TEXT,
                                    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                                    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                    PRIMARY KEY (`id`),
                                    KEY `idx_cv_collection` (`collection_id`),
                                    CONSTRAINT `fk_cv_collection`
                                        FOREIGN KEY (`collection_id`) REFERENCES `collections`(`id`)
                                            ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TABLE: diagram_collections
-- =========================================
CREATE TABLE `diagram_collections` (
                                       `id` CHAR(36) NOT NULL,
                                       `name` VARCHAR(255) NOT NULL,
                                       `description` TEXT,
                                       `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                                       `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                       `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TABLE: diagrams
-- =========================================
CREATE TABLE `diagrams` (
                            `id` CHAR(36) NOT NULL,
                            `collection_id` CHAR(36) NOT NULL,
                            `title` VARCHAR(255) NOT NULL,
                            `description` TEXT,
                            `syntax_type` VARCHAR(50) DEFAULT 'mermaid',
                            `syntax` TEXT NOT NULL,
                            `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                            `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `idx_diagrams_collection` (`collection_id`),
                            CONSTRAINT `fk_diagrams_collection`
                                FOREIGN KEY (`collection_id`) REFERENCES `collections`(`id`)
                                    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TABLE: endpoints
-- =========================================
CREATE TABLE `endpoints` (
                             `id` CHAR(36) NOT NULL,
                             `category` VARCHAR(75),
                             `collection_id` CHAR(36) NOT NULL,
                             `path` VARCHAR(75) NOT NULL,
                             `desc` VARCHAR(120),
                             `method` ENUM('GET','POST','PUT','PATCH','DELETE','OPTIONS','TRACE') NOT NULL,
                             `active_scenario` CHAR(36),
                             `enable_response_intercept` TINYINT(1) NOT NULL DEFAULT 0,
                             `delay` TINYINT,
                             `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             KEY `idx_endpoints_collection` (`collection_id`),
                             CONSTRAINT `fk_endpoints_collection`
                                 FOREIGN KEY (`collection_id`) REFERENCES `collections`(`id`)
                                     ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TABLE: scenario_response
-- =========================================
CREATE TABLE `scenario_response` (
                                     `id` CHAR(36) NOT NULL,
                                     `endpoint_id` CHAR(36) NOT NULL,
                                     `desc` VARCHAR(120),
                                     `status_header` INT,
                                     `delay` TINYINT,
                                     `body` TEXT,
                                     `header` LONGTEXT,
                                     `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
                                     `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                     PRIMARY KEY (`id`),
                                     KEY `idx_sr_endpoint` (`endpoint_id`),
                                     CONSTRAINT `fk_sr_endpoint`
                                         FOREIGN KEY (`endpoint_id`) REFERENCES `endpoints`(`id`)
                                             ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================
-- TRIGGERS (AUTO UUID)
-- =========================================
DELIMITER //

CREATE TRIGGER bi_squads
    BEFORE INSERT ON squads
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER bi_collections
    BEFORE INSERT ON collections
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER bi_custom_variables
    BEFORE INSERT ON custom_variables
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER bi_diagram_collections
    BEFORE INSERT ON diagram_collections
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER bi_diagrams
    BEFORE INSERT ON diagrams
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER bi_endpoints
    BEFORE INSERT ON endpoints
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

CREATE TRIGGER bi_scenario_response
    BEFORE INSERT ON scenario_response
    FOR EACH ROW
BEGIN
    IF NEW.id IS NULL OR NEW.id = '' THEN
        SET NEW.id = UUID();
END IF;
END//

DELIMITER ;