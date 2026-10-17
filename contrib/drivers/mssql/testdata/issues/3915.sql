CREATE TABLE [issue3915] (
    [id] int NOT NULL,
    [a] real NULL DEFAULT NULL,
    [b] real NULL DEFAULT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO [issue3915] ([id], [a], [b]) VALUES (1, 1, 2);
INSERT INTO [issue3915] ([id], [a], [b]) VALUES (2, 5, 4);
