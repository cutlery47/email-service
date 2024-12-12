CREATE TABLE IF NOT EXISTS email_schema.users (
    id              email_schema.uuid_key        PRIMARY KEY,
    mail            email_schema.string          NOT NULL,
    nickname        email_schema.string          NOT NULL,
    firstname       email_schema.string          NOT NULL,
    lastname        email_schema.string          NOT NULL,

    UNIQUE(mail)
);