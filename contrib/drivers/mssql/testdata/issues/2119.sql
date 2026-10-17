DROP TABLE IF EXISTS [sys_role];
CREATE TABLE [sys_role] (
    [id] int NOT NULL,
    [name] nvarchar(30) NOT NULL DEFAULT '',
    [code] nvarchar(100) NOT NULL DEFAULT '',
    [description] nvarchar(500) NOT NULL DEFAULT '',
    [weight] int NOT NULL DEFAULT 0,
    [status_id] int NOT NULL DEFAULT 1,
    [created_at] datetime2(0) NULL DEFAULT NULL,
    [updated_at] datetime2(0) NULL DEFAULT NULL,
    PRIMARY KEY ([id])
);
CREATE INDEX [code] ON [sys_role] ([code]);
INSERT INTO [sys_role] ([id], [name], [code], [description], [weight], [status_id], [created_at], [updated_at]) VALUES (1, N'开发人员', 'developer', '123123', 900, 2, '2022-09-03 21:25:03', '2022-09-09 23:35:23');
INSERT INTO [sys_role] ([id], [name], [code], [description], [weight], [status_id], [created_at], [updated_at]) VALUES (2, N'管理员', 'admin', '', 800, 1, '2022-09-03 21:25:03', '2022-09-09 23:00:17');
INSERT INTO [sys_role] ([id], [name], [code], [description], [weight], [status_id], [created_at], [updated_at]) VALUES (3, N'运营', 'operator', '', 700, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO [sys_role] ([id], [name], [code], [description], [weight], [status_id], [created_at], [updated_at]) VALUES (4, N'客服', 'service', '', 600, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');
INSERT INTO [sys_role] ([id], [name], [code], [description], [weight], [status_id], [created_at], [updated_at]) VALUES (5, N'收银', 'account', '', 500, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03');

DROP TABLE IF EXISTS [sys_status];
CREATE TABLE [sys_status] (
    [id] int NOT NULL,
    [en] nvarchar(50) NOT NULL DEFAULT '',
    [cn] nvarchar(50) NOT NULL DEFAULT '',
    [weight] int NOT NULL DEFAULT 0,
    PRIMARY KEY ([id])
);
INSERT INTO [sys_status] ([id], [en], [cn], [weight]) VALUES (1, 'on line', N'上线', 900);
INSERT INTO [sys_status] ([id], [en], [cn], [weight]) VALUES (2, 'undecided', N'未决定', 800);
INSERT INTO [sys_status] ([id], [en], [cn], [weight]) VALUES (3, 'off line', N'下线', 700);
