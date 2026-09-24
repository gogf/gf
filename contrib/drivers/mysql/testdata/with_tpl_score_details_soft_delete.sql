CREATE TABLE IF NOT EXISTS %s (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    score_id int(10) unsigned NOT NULL,
    detail_info varchar(100) NOT NULL,
    rank int(10) unsigned NOT NULL DEFAULT 0,
    deleted_at datetime default NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
