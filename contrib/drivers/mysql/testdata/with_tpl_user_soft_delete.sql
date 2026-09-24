CREATE TABLE IF NOT EXISTS %s (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    name varchar(45) NOT NULL,
    status int(10) unsigned NOT NULL DEFAULT 1,
    deleted_at datetime default NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
