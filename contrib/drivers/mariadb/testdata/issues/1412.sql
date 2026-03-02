-- ----------------------------
-- Table structure for items
-- ----------------------------
CREATE TABLE `items`  (
    `id` int(11) NOT NULL,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of items
-- ----------------------------
INSERT INTO `items` VALUES (1, '金秋产品1');
INSERT INTO `items` VALUES (2, '金秋产品2');

-- ----------------------------
-- Table structure for parcels
-- ----------------------------
CREATE TABLE `parcels`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `item_id` int(11) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of parcels
-- ----------------------------
INSERT INTO `parcels` VALUES (1, 1);
INSERT INTO `parcels` VALUES (2, 2);
INSERT INTO `parcels` VALUES (3, 0);