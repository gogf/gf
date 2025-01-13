DROP TABLE IF EXISTS `proxy_param`;
CREATE TABLE `proxy_param` (
                               `proxy_id` bigint NOT NULL,
                               `recommend_ids` json DEFAULT NULL,
                               `photos` json DEFAULT NULL,
                               PRIMARY KEY (`proxy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `proxy_param` (`proxy_id`, `recommend_ids`, `photos`) VALUES (1, '[584, 585]', 'null');
INSERT INTO `proxy_param` (`proxy_id`, `recommend_ids`, `photos`) VALUES (2, '[]', NULL);