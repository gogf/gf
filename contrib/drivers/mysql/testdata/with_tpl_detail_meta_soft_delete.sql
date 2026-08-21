CREATE TABLE IF NOT EXISTS %s (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    detail_id int(10) unsigned NOT NULL,
    meta_key varchar(50) NOT NULL,
    meta_value varchar(100) NOT NULL,
    sort_order int(10) unsigned NOT NULL DEFAULT 0,
    deleted_at datetime default NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
