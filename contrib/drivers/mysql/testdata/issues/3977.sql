DROP TABLE IF EXISTS `issue3977`;
CREATE TABLE `issue3977` (
    `id` bigint NOT NULL,
    `username` varchar(255) DEFAULT NULL,
    `balance` decimal(10,2) DEFAULT NULL,
    `state`  bool DEFAULT NULL,
    `age` int DEFAULT NULL,
    `create_at` datetime(0) DEFAULT NULL,
    `update_at` datetime(0) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

INSERT INTO `issue3977` VALUES (1, "username1", 1.01, 1, 18, "2020-01-01 00:00:00", "2020-01-01 00:00:00");
INSERT INTO `issue3977` VALUES (2, "", 0, 0, 0, "2020-01-02 00:00:00", "2020-01-02 00:00:00");
INSERT INTO `issue3977` VALUES (3, NULL, NULL, NULL, NULL, NULL, NULL);
