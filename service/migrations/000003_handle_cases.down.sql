ALTER TABLE
    "expert_cases"
    DROP CONSTRAINT "expert_cases_case_id_foreign";
ALTER TABLE
    "expert_cases"
    DROP CONSTRAINT "expert_cases_expert_id_foreign";

DROP TABLE expert_cases CASCADE;