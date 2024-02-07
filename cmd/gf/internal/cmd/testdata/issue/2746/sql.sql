CREATE TABLE %s (
     `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
     `nickname` varchar(45) NOT NULL COMMENT 'User Nickname',
     `tag` json NOT NULL,
     `info` longtext DEFAULT NULL,
     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

