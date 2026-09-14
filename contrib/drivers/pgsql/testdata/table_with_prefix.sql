DROP TABLE IF EXISTS instance;
CREATE TABLE instance (
    f_id SERIAL NOT NULL PRIMARY KEY,
    name varchar(255) DEFAULT ''
);
INSERT INTO instance VALUES (1, 'john');
