ALTER TABLE
    "cameras"
    DROP CONSTRAINT "cameras_camera_type_id_foreign";
ALTER TABLE
    "cases"
    DROP CONSTRAINT "cases_camera_id_foreign";
ALTER TABLE
    "cases"
    DROP CONSTRAINT "cases_violation_id_foreign";
ALTER TABLE
    "cases"
    DROP CONSTRAINT "cases_transport_id_foreign";
ALTER TABLE
    "transports"
    DROP CONSTRAINT "transports_person_id_foreign";

DROP TABLE cases CASCADE;

DROP TABLE camera_types CASCADE;

DROP TABLE cameras CASCADE;

DROP TABLE violations CASCADE;

DROP TABLE transports CASCADE;

DROP TABLE persons CASCADE;

