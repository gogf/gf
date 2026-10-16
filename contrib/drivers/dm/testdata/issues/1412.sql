-- ----------------------------
-- Table structure for items
-- ----------------------------
CREATE TABLE items  (
    id INT NOT NULL,
    name VARCHAR(255) NULL DEFAULT NULL,
    PRIMARY KEY (id)
);

-- ----------------------------
-- Records of items
-- ----------------------------
INSERT INTO items VALUES (1, '金秋产品1');
INSERT INTO items VALUES (2, '金秋产品2');

-- ----------------------------
-- Table structure for parcels
-- ----------------------------
CREATE TABLE parcels  (
    id INT NOT NULL AUTO_INCREMENT,
    item_id INT NULL DEFAULT NULL,
    PRIMARY KEY (id)
) AUTO_INCREMENT = 4;

-- ----------------------------
-- Records of parcels
-- ----------------------------
INSERT INTO parcels VALUES (1, 1);
INSERT INTO parcels VALUES (2, 2);
INSERT INTO parcels VALUES (3, 0);