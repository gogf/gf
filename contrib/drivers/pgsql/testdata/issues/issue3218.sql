CREATE TABLE "issue3218_sys_config" (
    "id" serial NOT NULL,
    "name" varchar(255) DEFAULT NULL,
    "value" text DEFAULT NULL,
    "created_at" timestamp DEFAULT NULL,
    "updated_at" timestamp DEFAULT NULL,
    PRIMARY KEY ("id"),
    UNIQUE ("name")
);

INSERT INTO "issue3218_sys_config" VALUES (49, 'site', '{"banned_ip":"22","filings":"2222","fixed_page":"","site_name":"22","version":"22"}', '2023-12-19 14:08:25', '2023-12-19 14:08:25');
