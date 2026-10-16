-- ----------------------------
-- Table structure for parcel_items
-- ----------------------------
DROP TABLE IF EXISTS parcel_items;
CREATE TABLE parcel_items  (
    id INT NOT NULL,
    parcel_id INT NULL DEFAULT NULL,
    name VARCHAR(255) NULL DEFAULT NULL,
    PRIMARY KEY (id)
);

-- ----------------------------
-- Records of parcel_items
-- ----------------------------
INSERT INTO parcel_items VALUES (1, 1, '新品');
INSERT INTO parcel_items VALUES (2, 3, '新品2');

-- ----------------------------
-- Table structure for parcels
-- ----------------------------
DROP TABLE IF EXISTS parcels;
CREATE TABLE parcels  (
    id INT NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (id)
) AUTO_INCREMENT = 4;

-- ----------------------------
-- Records of parcels
-- ----------------------------
INSERT INTO parcels VALUES (1);
INSERT INTO parcels VALUES (2);
INSERT INTO parcels VALUES (3);