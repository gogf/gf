CREATE TABLE "issue3086_user" (
    "id" int NOT NULL,
    "passport" varchar(45) NOT NULL,
    "password" varchar(45) DEFAULT NULL,
    "nickname" varchar(45) DEFAULT NULL,
    "create_at" timestamp DEFAULT NULL,
    "update_at" timestamp DEFAULT NULL,
    PRIMARY KEY ("id")
);
