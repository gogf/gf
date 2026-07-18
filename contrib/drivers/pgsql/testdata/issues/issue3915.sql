CREATE TABLE "issue3915" (
    "id" serial NOT NULL,
    "a" real DEFAULT NULL,
    "b" real DEFAULT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "issue3915" ("id", "a", "b") VALUES (1, 1, 2);
INSERT INTO "issue3915" ("id", "a", "b") VALUES (2, 5, 4);
