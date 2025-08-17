CREATE TABLE `%s` (
    `Id`        int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
    `parentId`  varchar(45) NOT NULL COMMENT '',
    `PASSPORT`  varchar(45) NOT NULL COMMENT 'User Passport',
    `PASS_WORD`  varchar(45) NOT NULL COMMENT 'User Password',
    `NICKNAME2`  varchar(45) NOT NULL COMMENT 'User Nickname',
    `create_at` datetime DEFAULT NULL COMMENT 'Created Time',
    `update_at` datetime DEFAULT NULL COMMENT 'Updated Time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
