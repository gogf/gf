CREATE TABLE `date_time_example` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `year` year DEFAULT NULL COMMENT 'year',
    `date` date DEFAULT NULL COMMENT 'Date',
    `time` time DEFAULT NULL COMMENT 'time',
    `datetime` datetime DEFAULT NULL COMMENT 'datetime',
    `timestamp` timestamp NULL DEFAULT NULL COMMENT 'Timestamp',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;