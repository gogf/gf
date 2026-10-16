CREATE TABLE date_time_example (
    id INT NOT NULL AUTO_INCREMENT,
    year INT DEFAULT NULL COMMENT 'year',
    date DATE DEFAULT NULL COMMENT 'Date',
    time TIME DEFAULT NULL COMMENT 'time',
    datetime DATETIME DEFAULT NULL COMMENT 'datetime',
    timestamp TIMESTAMP NULL DEFAULT NULL COMMENT 'Timestamp',
    PRIMARY KEY (id)
);
