CREATE TABLE "issue2105" (
    "id" varchar(255) NOT NULL,
    "json" text,
    PRIMARY KEY ("id")
);

INSERT INTO "issue2105" VALUES ('1', NULL);
INSERT INTO "issue2105" VALUES ('2', '[{"Name": "任务类型", "Value": "高价值"}, {"Name": "优先级", "Value": "高"}, {"Name": "是否亮点功能", "Value": "是"}]');
