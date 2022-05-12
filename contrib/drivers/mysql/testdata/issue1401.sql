-- ----------------------------
-- Table structure for parcel_items
-- ----------------------------
DROP TABLE IF EXISTS `parcel_items`;
CREATE TABLE `parcel_items`  (
    `id` int(11) NOT NULL,
    `parcel_id` int(11) NULL DEFAULT NULL,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of parcel_items
-- ----------------------------
INSERT INTO `parcel_items` VALUES (1, 1, '新品');
INSERT INTO `parcel_items` VALUES (2, 3, '新品2');

-- ----------------------------
-- Table structure for parcels
-- ----------------------------
DROP TABLE IF EXISTS `parcels`;
CREATE TABLE `parcels`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of parcels
-- ----------------------------
INSERT INTO `parcels` VALUES (1);
INSERT INTO `parcels` VALUES (2);
INSERT INTO `parcels` VALUES (3);