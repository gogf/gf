DROP TABLE IF EXISTS `issue3977`;
CREATE TABLE `issue3977` (
    `id` bigint NOT NULL,
    `username` varchar(255) DEFAULT "",
    `balance` decimal(10,2) DEFAULT 0.00,
    `state`  bool DEFAULT 0,
    `age` int DEFAULT 0,
    `create_at` datetime(0) DEFAULT NULL,
    `update_at` datetime(0) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

INSERT INTO `issue3977` VALUES (1, "username1", 1.01, 1, 18, "2020-01-01 00:00:00", "2020-01-01 00:00:00");
INSERT INTO `issue3977` VALUES (2, "username2", 2.02, 1, 100, "2020-01-01 00:00:00", "2020-01-01 00:00:00");
INSERT INTO `issue3977` VALUES (3, "username3", 3.03, 0, 56, "2020-01-01 00:00:00", "2020-01-01 00:00:00");