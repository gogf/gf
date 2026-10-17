CREATE TABLE a (
    [id] int NOT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO a ([id]) VALUES (2);

CREATE TABLE b (
    [id] int NOT NULL,
    [name] nvarchar(255) NOT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO b ([id], [name]) VALUES (2, 'a');
INSERT INTO b ([id], [name]) VALUES (3, 'b');

CREATE TABLE c (
    [id] int NOT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO c ([id]) VALUES (2);
