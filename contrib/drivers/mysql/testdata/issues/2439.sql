CREATE TABLE `a`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (id) USING BTREE
) ENGINE = InnoDB;
INSERT INTO `a` (`id`) VALUES ('2');

CREATE TABLE `b`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL ,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB;
INSERT INTO `b` (`id`, `name`) VALUES ('2', 'a');
INSERT INTO `b` (`id`, `name`) VALUES ('3', 'b');

CREATE TABLE `c`  (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB;
INSERT INTO `c` (`id`) VALUES ('2');