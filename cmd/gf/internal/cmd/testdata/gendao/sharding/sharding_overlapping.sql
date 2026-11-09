-- Test case for issue #4603: overlapping sharding patterns
-- https://github.com/gogf/gf/issues/4603
--
-- Patterns: "a_?", "a_b_?", "a_c_?"
-- Expected: a_1/a_2 -> "a", a_b_1/a_b_2 -> "a_b", a_c_1/a_c_2 -> "a_c"

CREATE TABLE `a_1`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(45) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `a_2`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(45) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `a_b_1`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(45) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `a_b_2`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(45) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `a_c_1`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(45) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `a_c_2`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(45) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
