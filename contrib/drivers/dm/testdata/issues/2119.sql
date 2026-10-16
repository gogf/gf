-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS sys_role;
CREATE TABLE sys_role  (
                             id INT NOT NULL AUTO_INCREMENT COMMENT '||s',
                             name VARCHAR(30) NOT NULL DEFAULT '' COMMENT '角色名称||s,r',
                             code VARCHAR(100) NOT NULL DEFAULT '' COMMENT '角色 code||s,r',
                             description VARCHAR(500) NOT NULL DEFAULT '' COMMENT '描述信息|text',
                             weight INT NOT NULL DEFAULT 0 COMMENT '排序||r|min:0#发布状态不能小于 0',
                             status_id INT NOT NULL DEFAULT 1 COMMENT '发布状态|hasOne|f:status,fk:id',
                             created_at DATETIME(0) NULL DEFAULT NULL,
                             updated_at DATETIME(0) NULL DEFAULT NULL,
                             PRIMARY KEY (id)
) AUTO_INCREMENT = 1091;
CREATE INDEX code ON sys_role (code);
COMMENT ON TABLE sys_role IS '系统角色表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO sys_role VALUES (1, '开发人员', 'developer', '123123', 900, 2, '2022-09-03 21:25:03', '2022-09-09 23:35:23');
INSERT INTO sys_role VALUES (2, '管理员', 'admin', '', 800, 1, '2022-09-03 21:25:03', '2022-09-09 23:00:17');
INSERT INTO sys_role VALUES (3, '运营', 'operator', '', 700, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO sys_role VALUES (4, '客服', 'service', '', 600, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO sys_role VALUES (5, '收银', 'account', '', 500, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');

-- ----------------------------
-- Table structure for sys_status
-- ----------------------------
DROP TABLE IF EXISTS sys_status;
CREATE TABLE sys_status  (
                               id INT NOT NULL AUTO_INCREMENT,
                               en VARCHAR(50) NOT NULL DEFAULT '' COMMENT '英文名称',
                               cn VARCHAR(50) NOT NULL DEFAULT '' COMMENT '中文名称',
                               weight INT NOT NULL DEFAULT 0 COMMENT '排序权重',
                               PRIMARY KEY (id)
) AUTO_INCREMENT = 7;
COMMENT ON TABLE sys_status IS '发布状态';

-- ----------------------------
-- Records of sys_status
-- ----------------------------
INSERT INTO sys_status VALUES (1, 'on line', '上线', 900);
INSERT INTO sys_status VALUES (2, 'undecided', '未决定', 800);
INSERT INTO sys_status VALUES (3, 'off line', '下线', 700);