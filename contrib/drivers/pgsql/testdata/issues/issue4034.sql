CREATE TABLE "issue4034" (
    "id" serial NOT NULL,
    "passport" varchar(255),
    "password" varchar(255),
    "nickname" varchar(255),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);
