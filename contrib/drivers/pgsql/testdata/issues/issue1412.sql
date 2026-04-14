CREATE TABLE "items" (
    "id" int NOT NULL,
    "name" varchar(255) DEFAULT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "items" VALUES (1, '金秋产品1');
INSERT INTO "items" VALUES (2, '金秋产品2');

CREATE TABLE "parcels" (
    "id" serial NOT NULL,
    "item_id" int DEFAULT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "parcels" VALUES (1, 1);
INSERT INTO "parcels" VALUES (2, 2);
INSERT INTO "parcels" VALUES (3, 0);
