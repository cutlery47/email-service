CREATE SCHEMA IF NOT EXISTS email_schema;

CREATE DOMAIN email_schema.uuid_key AS UUID
DEFAULT gen_random_uuid()
NOT NULL;

CREATE DOMAIN email_schema.string AS VARCHAR(256)
NOT NULL;