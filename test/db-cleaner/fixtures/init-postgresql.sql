CREATE SCHEMA "schema" AUTHORIZATION postgres;

CREATE TABLE "schema".table_a (id BIGINT PRIMARY KEY, text VARCHAR(255) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE TABLE "schema".table_b (id BIGINT PRIMARY KEY, text VARCHAR(255) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE TABLE "schema".table_c (id BIGINT PRIMARY KEY, text VARCHAR(255) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now());

INSERT INTO "schema".table_a(id, text, created_at) VALUES (1, 'text', now() - INTERVAL '1 days');
INSERT INTO "schema".table_a(id, text, created_at) VALUES (2, 'text', now() - INTERVAL '2 days');
INSERT INTO "schema".table_a(id, text, created_at) VALUES (3, 'text', now() - INTERVAL '3 days');

INSERT INTO "schema".table_b(id, text, created_at) VALUES (1, 'text', now() - INTERVAL '1 days');
INSERT INTO "schema".table_b(id, text, created_at) VALUES (2, 'text', now() - INTERVAL '2 days');
INSERT INTO "schema".table_b(id, text, created_at) VALUES (3, 'text', now() - INTERVAL '3 days');

INSERT INTO "schema".table_c(id, text, created_at) VALUES (1, 'text', now() - INTERVAL '1 days');
INSERT INTO "schema".table_c(id, text, created_at) VALUES (2, 'text', now() - INTERVAL '2 days');
INSERT INTO "schema".table_c(id, text, created_at) VALUES (3, 'text', now() - INTERVAL '3 days');
