CREATE TABLE sys_role (
    id          NUMBER(10)     NOT NULL,
    name        VARCHAR2(30),
    code        VARCHAR2(100),
    description VARCHAR2(500),
    weight      NUMBER(10)     DEFAULT 0 NOT NULL,
    status_id   NUMBER(10)     DEFAULT 1 NOT NULL,
    created_at  TIMESTAMP      DEFAULT NULL,
    updated_at  TIMESTAMP      DEFAULT NULL,
    PRIMARY KEY (id)
);
INSERT INTO sys_role VALUES (1, '开发人员', 'developer', '123123', 900, 2, TIMESTAMP '2022-09-03 21:25:03', TIMESTAMP '2022-09-09 23:35:23');
INSERT INTO sys_role VALUES (2, '管理员', 'admin', '', 800, 1, TIMESTAMP '2022-09-03 21:25:03', TIMESTAMP '2022-09-09 23:00:17');
INSERT INTO sys_role VALUES (3, '运营', 'operator', '', 700, 1, TIMESTAMP '2022-09-03 21:25:03', TIMESTAMP '2022-09-03 21:25:03');
INSERT INTO sys_role VALUES (4, '客服', 'service', '', 600, 1, TIMESTAMP '2022-09-03 21:25:03', TIMESTAMP '2022-09-03 21:25:03');
INSERT INTO sys_role VALUES (5, '收银', 'account', '', 500, 1, TIMESTAMP '2022-09-03 21:25:03', TIMESTAMP '2022-09-03 21:25:03');
CREATE TABLE sys_status (
    id     NUMBER(10)    NOT NULL,
    en     VARCHAR2(50),
    cn     VARCHAR2(50),
    weight NUMBER(10)    DEFAULT 0 NOT NULL,
    PRIMARY KEY (id)
);
INSERT INTO sys_status VALUES (1, 'on line', '上线', 900);
INSERT INTO sys_status VALUES (2, 'undecided', '未决定', 800);
INSERT INTO sys_status VALUES (3, 'off line', '下线', 700);
