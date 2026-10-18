CREATE TABLE parcel_items (
    id        NUMBER(10) NOT NULL,
    parcel_id NUMBER(10) DEFAULT NULL,
    name      VARCHAR2(255) DEFAULT NULL,
    PRIMARY KEY (id)
);
INSERT INTO parcel_items VALUES (1, 1, '新品');
INSERT INTO parcel_items VALUES (2, 3, '新品2');
CREATE TABLE parcels (
    id NUMBER(10) NOT NULL,
    PRIMARY KEY (id)
);
INSERT INTO parcels VALUES (1);
INSERT INTO parcels VALUES (2);
INSERT INTO parcels VALUES (3);
