-- DDL Script for Gigawrks Database --

CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `country` varchar(50) NOT NULL,
  `email` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email_UNIQUE` (`email`)
);

CREATE TABLE IF NOT EXISTS `countries`(
  `id` int NOT NULL AUTO_INCREMENT,
  `common_name` varchar(50) NOT NULL,
  `official_name` varchar(100) DEFAULT NULL,
  `country_code` varchar(30) NOT NULL,
  `capital` varchar(50) DEFAULT NULL,
  `region` varchar(50) DEFAULT NULL,
  `sub_region` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `country_code_UNIQUE` (`country_code`)
)