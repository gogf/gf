-- ----------------------------
-- Table structure for issue2012
-- ----------------------------
DROP TABLE IF EXISTS `issue2012`;

CREATE TABLE `issue2012`(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `time_only` time,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
