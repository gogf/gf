
CREATE TABLE `table_a`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `alias` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB;

INSERT INTO `table_a` VALUES (1, 'table_a_test1');
INSERT INTO `table_a` VALUES (2, 'table_a_test2');

CREATE TABLE `table_b`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `table_a_id` int(11) NOT NULL,
    `alias` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB;

INSERT INTO `table_b` VALUES (10, 1, 'table_b_test1');
INSERT INTO `table_b` VALUES (20, 2, 'table_b_test2');
INSERT INTO `table_b` VALUES (30, 1, 'table_b_test3');
INSERT INTO `table_b` VALUES (40, 2, 'table_b_test4');

CREATE TABLE `table_c`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `table_b_id` int(11) NOT NULL,
    `alias` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB;

INSERT INTO `table_c` VALUES (100, 10, 'table_c_test1');
INSERT INTO `table_c` VALUES (200, 10, 'table_c_test2');
INSERT INTO `table_c` VALUES (300, 20, 'table_c_test3');
INSERT INTO `table_c` VALUES (400, 30, 'table_c_test4');