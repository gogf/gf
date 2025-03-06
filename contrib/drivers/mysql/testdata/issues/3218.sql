CREATE TABLE `issue3218_sys_config`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '配置名称',
    `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '配置值',
    `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `name`(`name`(191)) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `issue3218_sys_config` VALUES (49, 'site', '{\"banned_ip\":\"22\",\"filings\":\"2222\",\"fixed_page\":\"\",\"site_name\":\"22\",\"version\":\"22\"}', '2023-12-19 14:08:25', '2023-12-19 14:08:25');