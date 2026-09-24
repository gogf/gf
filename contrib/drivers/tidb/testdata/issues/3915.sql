CREATE TABLE `issue3915` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'user id',
    `a` float DEFAULT NULL COMMENT 'user name',
    `b` float DEFAULT NULL COMMENT 'user status',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

INSERT INTO `issue3915` (`id`,`a`,`b`) VALUES (1,1,2);
INSERT INTO `issue3915` (`id`,`a`,`b`) VALUES (2,5,4);