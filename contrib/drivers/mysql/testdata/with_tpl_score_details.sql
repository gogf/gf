CREATE TABLE IF NOT EXISTS %s (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    score_id int(10) unsigned NOT NULL,
    detail_info varchar(45) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
