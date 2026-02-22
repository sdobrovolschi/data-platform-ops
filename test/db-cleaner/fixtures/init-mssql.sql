CREATE DATABASE [database];
GO
USE [database];
GO
CREATE SCHEMA [schema] AUTHORIZATION dbo;
GO

CREATE TABLE [schema].table_a (id BIGINT PRIMARY KEY, text VARCHAR(255) NOT NULL, created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME());
CREATE TABLE [schema].table_b (id BIGINT PRIMARY KEY, text VARCHAR(255) NOT NULL, created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME());
CREATE TABLE [schema].table_c (id BIGINT PRIMARY KEY, text VARCHAR(255) NOT NULL, created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME());

INSERT INTO [schema].table_a(id, text, created_at) VALUES (1, 'text', DATEADD(DAY, -1, SYSDATETIME()));
INSERT INTO [schema].table_a(id, text, created_at) VALUES (2, 'text', DATEADD(DAY, -2, SYSDATETIME()));
INSERT INTO [schema].table_a(id, text, created_at) VALUES (3, 'text', DATEADD(DAY, -3, SYSDATETIME()));

INSERT INTO [schema].table_b(id, text, created_at) VALUES (1, 'text', DATEADD(DAY, -1, SYSDATETIME()));
INSERT INTO [schema].table_b(id, text, created_at) VALUES (2, 'text', DATEADD(DAY, -2, SYSDATETIME()));
INSERT INTO [schema].table_b(id, text, created_at) VALUES (3, 'text', DATEADD(DAY, -3, SYSDATETIME()));

INSERT INTO [schema].table_c(id, text, created_at) VALUES (1, 'text', DATEADD(DAY, -1, SYSDATETIME()));
INSERT INTO [schema].table_c(id, text, created_at) VALUES (2, 'text', DATEADD(DAY, -2, SYSDATETIME()));
INSERT INTO [schema].table_c(id, text, created_at) VALUES (3, 'text', DATEADD(DAY, -3, SYSDATETIME()));
