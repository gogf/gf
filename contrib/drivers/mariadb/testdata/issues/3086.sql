CREATE TABLE `issue3086_user`
(
    `id`        int(10) unsigned NOT NULL COMMENT 'User ID',
    `passport`  varchar(45) NOT NULL COMMENT 'User Passport',
    `password`  varchar(45) DEFAULT NULL COMMENT 'User Password',
    `nickname`  varchar(45) DEFAULT NULL COMMENT 'User Nickname',
    `create_at` datetime DEFAULT NULL COMMENT 'Created Time',
    `update_at` datetime DEFAULT NULL COMMENT 'Updated Time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
