
CREATE TABLE table_a (
    id SERIAL PRIMARY KEY,
    alias varchar(255) DEFAULT ''
);

INSERT INTO table_a VALUES (1, 'table_a_test1');
INSERT INTO table_a VALUES (2, 'table_a_test2');

CREATE TABLE table_b (
    id SERIAL PRIMARY KEY,
    table_a_id integer NOT NULL,
    alias varchar(255) DEFAULT ''
);

INSERT INTO table_b VALUES (10, 1, 'table_b_test1');
INSERT INTO table_b VALUES (20, 2, 'table_b_test2');
INSERT INTO table_b VALUES (30, 1, 'table_b_test3');
INSERT INTO table_b VALUES (40, 2, 'table_b_test4');

CREATE TABLE table_c (
    id SERIAL PRIMARY KEY,
    table_b_id integer NOT NULL,
    alias varchar(255) DEFAULT ''
);

INSERT INTO table_c VALUES (100, 10, 'table_c_test1');
INSERT INTO table_c VALUES (200, 10, 'table_c_test2');
INSERT INTO table_c VALUES (300, 20, 'table_c_test3');
INSERT INTO table_c VALUES (400, 30, 'table_c_test4');
