CREATE TABLE IF NOT EXISTS %s (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    uid int(10) unsigned NOT NULL,
    score int(10) unsigned NOT NULL,
    priority int(10) unsigned NOT NULL DEFAULT 0,
    deleted_at datetime default NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
