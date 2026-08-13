DROP TABLE IF EXISTS "parcel_items";
CREATE TABLE "parcel_items" (
    "id" int NOT NULL,
    "parcel_id" int DEFAULT NULL,
    "name" varchar(255) DEFAULT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "parcel_items" VALUES (1, 1, '新品');
INSERT INTO "parcel_items" VALUES (2, 3, '新品2');

DROP TABLE IF EXISTS "parcels";
CREATE TABLE "parcels" (
    "id" serial NOT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "parcels" VALUES (1);
INSERT INTO "parcels" VALUES (2);
INSERT INTO "parcels" VALUES (3);
