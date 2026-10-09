CREATE TABLE `instance`  (
    `f_id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NULL DEFAULT '',
    PRIMARY KEY (`f_id`) USING BTREE
) ENGINE = InnoDB;

INSERT INTO `instance` VALUES (1, 'john');