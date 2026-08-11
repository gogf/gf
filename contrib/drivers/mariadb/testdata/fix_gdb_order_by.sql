CREATE TABLE IF NOT EXISTS `employee`
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255)                        NOT NULL,
    age        INT                                 NOT NULL
);

INSERT INTO employee(name, age) VALUES ('John', 30);
INSERT INTO employee(name, age) VALUES ('Mary', 28);