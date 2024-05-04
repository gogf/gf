SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
                             `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '||s',
                             `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色名称||s,r',
                             `code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色 code||s,r',
                             `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述信息|text',
                             `weight` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序||r|min:0#发布状态不能小于 0',
                             `status_id` int(0) UNSIGNED NOT NULL DEFAULT 1 COMMENT '发布状态|hasOne|f:status,fk:id',
                             `created_at` datetime(0) NULL DEFAULT NULL,
                             `updated_at` datetime(0) NULL DEFAULT NULL,
                             PRIMARY KEY (`id`) USING BTREE,
                             INDEX `code`(`code`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1091 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统角色表' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, '开发人员', 'developer', '123123', 900, 2, '2022-09-03 21:25:03', '2022-09-09 23:35:23');
INSERT INTO `sys_role` VALUES (2, '管理员', 'admin', '', 800, 1, '2022-09-03 21:25:03', '2022-09-09 23:00:17');
INSERT INTO `sys_role` VALUES (3, '运营', 'operator', '', 700, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO `sys_role` VALUES (4, '客服', 'service', '', 600, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO `sys_role` VALUES (5, '收银', 'account', '', 500, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');

-- ----------------------------
-- Table structure for sys_status
-- ----------------------------
DROP TABLE IF EXISTS `sys_status`;
CREATE TABLE `sys_status`  (
                               `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
                               `en` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '英文名称',
                               `cn` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '中文名称',
                               `weight` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序权重',
                               PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '发布状态' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_status
-- ----------------------------
INSERT INTO `sys_status` VALUES (1, 'on line', '上线', 900);
INSERT INTO `sys_status` VALUES (2, 'undecided', '未决定', 800);
INSERT INTO `sys_status` VALUES (3, 'off line', '下线', 700);