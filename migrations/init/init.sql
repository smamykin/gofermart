CREATE TABLE "user" (
    "id" SERIAL PRIMARY KEY,
    "login" VARCHAR NOT NULL ,
    "pwd" VARCHAR NOT NULL
);

CREATE UNIQUE INDEX name_type_unique ON "user" (login);