-- Migration: add `script` column to endpoints
--
-- Adds the per-endpoint JavaScript snippet used by the API-Mock script feature
-- (executed via goja on every mock request). Run this against existing databases;
-- fresh installs already include the column via database.sql / mariadb.sql.
--
-- MySQL 8 (does not support ADD COLUMN IF NOT EXISTS):
ALTER TABLE `endpoints` ADD COLUMN `script` TEXT NULL AFTER `active_scenario`;

-- MariaDB equivalent (idempotent):
-- ALTER TABLE `endpoints` ADD COLUMN IF NOT EXISTS `script` TEXT NULL AFTER `active_scenario`;
