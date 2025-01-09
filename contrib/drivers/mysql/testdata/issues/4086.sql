DROP TABLE IF EXISTS `issue4086`;
CREATE TABLE `issue4086` (
    `proxy_id` bigint NOT NULL,
    `recommend_ids` json DEFAULT NULL,
    `photos` json DEFAULT NULL,
    PRIMARY KEY (`proxy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `issue4086` (`proxy_id`, `recommend_ids`, `photos`) VALUES (1, '[584, 585]', 'null');
INSERT INTO `issue4086` (`proxy_id`, `recommend_ids`, `photos`) VALUES (2, '[]', NULL);