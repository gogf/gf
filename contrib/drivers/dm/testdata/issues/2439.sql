CREATE TABLE a (
    id NUMBER(10) NOT NULL,
    PRIMARY KEY (id)
);
INSERT INTO a (id) VALUES ('2');
CREATE TABLE b (
    id   NUMBER(10) NOT NULL,
    name VARCHAR2(255) NOT NULL,
    PRIMARY KEY (id)
);
INSERT INTO b (id, name) VALUES ('2', 'a');
INSERT INTO b (id, name) VALUES ('3', 'b');
CREATE TABLE c (
    id NUMBER(10) NOT NULL,
    PRIMARY KEY (id)
);
INSERT INTO c (id) VALUES ('2');
