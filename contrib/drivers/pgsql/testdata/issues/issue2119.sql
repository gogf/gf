DROP TABLE IF EXISTS "sys_role";
CREATE TABLE "sys_role" (
    "id" serial NOT NULL,
    "name" varchar(30) NOT NULL DEFAULT '',
    "code" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "weight" int NOT NULL DEFAULT 0,
    "status_id" int NOT NULL DEFAULT 1,
    "created_at" timestamp DEFAULT NULL,
    "updated_at" timestamp DEFAULT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "sys_role" VALUES (1, '开发人员', 'developer', '123123', 900, 2, '2022-09-03 21:25:03', '2022-09-09 23:35:23');
INSERT INTO "sys_role" VALUES (2, '管理员', 'admin', '', 800, 1, '2022-09-03 21:25:03', '2022-09-09 23:00:17');
INSERT INTO "sys_role" VALUES (3, '运营', 'operator', '', 700, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO "sys_role" VALUES (4, '客服', 'service', '', 600, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO "sys_role" VALUES (5, '收银', 'account', '', 500, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');

DROP TABLE IF EXISTS "sys_status";
CREATE TABLE "sys_status" (
    "id" serial NOT NULL,
    "en" varchar(50) NOT NULL DEFAULT '',
    "cn" varchar(50) NOT NULL DEFAULT '',
    "weight" int NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

INSERT INTO "sys_status" VALUES (1, 'on line', '上线', 900);
INSERT INTO "sys_status" VALUES (2, 'undecided', '未决定', 800);
INSERT INTO "sys_status" VALUES (3, 'off line', '下线', 700);
