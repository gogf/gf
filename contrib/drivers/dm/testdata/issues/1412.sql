CREATE TABLE items (
    id   NUMBER(10) NOT NULL,
    name VARCHAR2(255) DEFAULT NULL,
    PRIMARY KEY (id)
);
INSERT INTO items VALUES (1, '金秋产品1');
INSERT INTO items VALUES (2, '金秋产品2');
CREATE TABLE parcels (
    id      NUMBER(10) NOT NULL,
    item_id NUMBER(10) DEFAULT NULL,
    PRIMARY KEY (id)
);
INSERT INTO parcels VALUES (1, 1);
INSERT INTO parcels VALUES (2, 2);
INSERT INTO parcels VALUES (3, 0);
