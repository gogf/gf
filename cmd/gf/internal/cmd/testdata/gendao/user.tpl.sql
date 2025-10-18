CREATE TABLE `%s` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
    `passport` varchar(45) NOT NULL COMMENT 'User Passport',
    `password` varchar(45) NOT NULL COMMENT 'User Password',
    `nickname` varchar(45) NOT NULL COMMENT 'User Nickname',
    `score` decimal(10,2) unsigned DEFAULT NULL COMMENT 'Total score amount.',
    `create_at` datetime DEFAULT NULL COMMENT 'Created Time',
    `update_at` datetime DEFAULT NULL COMMENT 'Updated Time',
    `email` varchar(255) DEFAULT NULL COMMENT 'User Email',
    `status` int DEFAULT NULL COMMENT 'User Status',
    `height` float DEFAULT NULL COMMENT 'User Height',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
