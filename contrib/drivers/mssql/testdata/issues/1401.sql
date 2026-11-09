DROP TABLE IF EXISTS [parcel_items];
CREATE TABLE [parcel_items] (
    [id] int NOT NULL,
    [parcel_id] int NULL DEFAULT NULL,
    [name] nvarchar(255) NULL DEFAULT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO [parcel_items] ([id], [parcel_id], [name]) VALUES (1, 1, N'新品');
INSERT INTO [parcel_items] ([id], [parcel_id], [name]) VALUES (2, 3, N'新品2');

DROP TABLE IF EXISTS [parcels];
CREATE TABLE [parcels] (
    [id] int NOT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO [parcels] ([id]) VALUES (1);
INSERT INTO [parcels] ([id]) VALUES (2);
INSERT INTO [parcels] ([id]) VALUES (3);
