CREATE TABLE issue3218_sys_config  (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NULL DEFAULT NULL COMMENT '配置名称',
    value TEXT NULL COMMENT '配置值',
    created_at TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
    updated_at TIMESTAMP NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE (name)
);

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO issue3218_sys_config VALUES (49, 'site', '{"banned_ip":"22","filings":"2222","fixed_page":"","site_name":"22","version":"22"}', '2023-12-19 14:08:25', '2023-12-19 14:08:25');