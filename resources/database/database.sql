-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: mysql
-- Generation Time: Feb 03, 2023 at 02:38 AM
-- Server version: 8.0.32
-- PHP Version: 8.0.27
SET
    SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

START TRANSACTION;

SET
    time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */
;

/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */
;

/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */
;

/*!40101 SET NAMES utf8mb4 */
;

--
-- Database: `bonclay`
--
-- --------------------------------------------------------
--
-- Table structure for table `collections`
--
CREATE TABLE `collections` (
    `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `squad_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(75) COLLATE utf8mb4_unicode_ci NOT NULL,
    `baseURL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `slug` varchar(125) COLLATE utf8mb4_unicode_ci NOT NULL,
    `desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `docs` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- --------------------------------------------------------
--
-- Table structure for table `endpoints`
--
CREATE TABLE `endpoints` (
    `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `category` varchar(75) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `collection_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `path` varchar(75) COLLATE utf8mb4_unicode_ci NOT NULL,
    `desc` varchar(120) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `method` enum(
        'GET',
        'POST',
        'PUT',
        'PATCH',
        'DELETE',
        'OPTIONS',
        'TRACE'
    ) COLLATE utf8mb4_unicode_ci NOT NULL,
    `active_scenario` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- --------------------------------------------------------
--
-- Table structure for table `migrations`
--
CREATE TABLE `migrations` (
    `id` int UNSIGNED NOT NULL,
    `migration` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `batch` int NOT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- --------------------------------------------------------
--
-- Table structure for table `scenario_response`
--
CREATE TABLE `scenario_response` (
    `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `endpoint_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `desc` varchar(120) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `status_header` int DEFAULT NULL,
    `delay` tinyint DEFAULT NULL,
    `body` text COLLATE utf8mb4_unicode_ci,
    `header` json DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- --------------------------------------------------------
--
-- Table structure for table `squads`
--
CREATE TABLE `squads` (
    `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(75) COLLATE utf8mb4_unicode_ci NOT NULL,
    `slug` varchar(125) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

--
-- Indexes for dumped tables
--
--
-- Indexes for table `migrations`
--
ALTER TABLE
    `migrations`
ADD
    PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--
--
-- AUTO_INCREMENT for table `migrations`
--
ALTER TABLE
    `migrations`
MODIFY
    `id` int UNSIGNED NOT NULL AUTO_INCREMENT;

COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */
;

/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */
;

/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */
;

CREATE TABLE bonclay.custom_variables (
    id char(36) NOT NULL,
    collection_id char(36) NULL,
    `key` varchar(255) NOT NULL,
    value TEXT NULL,
    created_at timestamp NULL,
    deleted_at timestamp NULL,
    updated_at timestamp NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;


ALTER TABLE bonclay.endpoints ADD script TEXT NULL;
ALTER TABLE bonclay.endpoints CHANGE script script TEXT NULL AFTER active_scenario;

ALTER TABLE bonclay.custom_variables ADD CONSTRAINT custom_variables_pk PRIMARY KEY (id);
ALTER TABLE bonclay.custom_variables ADD CONSTRAINT custom_variables_unique UNIQUE KEY (collection_id,`key`);
