CREATE TABLE "users"
(
    "user_id"     UUID                               NOT NULL,
    "username"    VARCHAR(255)                       NOT NULL,
    "hash_pass"   VARCHAR(255)                       NOT NULL,
    "register_at" TIME(0) WITHOUT TIME ZONE          NOT NULL DEFAULT current_time,
    "role"        VARCHAR(255) CHECK
        ("role" IN ('director', 'expert', 'camera')) NOT NULL
);
ALTER TABLE
    "users"
    ADD PRIMARY KEY ("user_id");
CREATE TABLE "directors"
(
    "director_id" UUID NOT NULL,
    "user_id"     UUID NOT NULL
);
ALTER TABLE
    "directors"
    ADD PRIMARY KEY ("director_id");

CREATE TABLE "experts"
(
    "expert_id"        UUID    NOT NULL,
    "is_confirmed"     BOOLEAN NOT NULL,
    "user_id"          UUID    NOT NULL,
    "competence_skill" BIGINT  NOT NULL
);
ALTER TABLE
    "experts"
    ADD PRIMARY KEY ("expert_id");

ALTER TABLE
    "directors"
    ADD CONSTRAINT "directors_user_id_foreign" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
ALTER TABLE
    "experts"
    ADD CONSTRAINT "experts_user_id_foreign" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE cameras
    ADD COLUMN user_id UUID;

ALTER TABLE
    "cameras"
    ADD CONSTRAINT "cameras_user_id_foreign" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");