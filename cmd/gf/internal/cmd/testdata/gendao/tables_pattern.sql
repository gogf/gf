-- Test case for issue #4629: tables pattern matching
-- https://github.com/gogf/gf/issues/4629
-- Standard SQL syntax compatible with MySQL and PostgreSQL
--
-- Tables: trade_order, trade_item, user_info, user_log, config

CREATE TABLE trade_order (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(45) NOT NULL
);

CREATE TABLE trade_item (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(45) NOT NULL
);

CREATE TABLE user_info (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(45) NOT NULL
);

CREATE TABLE user_log (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(45) NOT NULL
);

CREATE TABLE config (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(45) NOT NULL
);
