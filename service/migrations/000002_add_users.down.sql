ALTER TABLE
    "directors"
    DROP CONSTRAINT "directors_user_id_foreign";
ALTER TABLE
    "experts"
    DROP CONSTRAINT "experts_user_id_foreign";

ALTER TABLE
    "cameras"
    DROP CONSTRAINT "cameras_user_id_foreign";

DROP TABLE users CASCADE;

DROP TABLE directors CASCADE;

DROP TABLE experts CASCADE;