CREATE TABLE "expert_cases"
(
    "expert_case_id"  UUID    NOT NULL,
    "expert_id"       UUID    NOT NULL,
    "case_id"         UUID    NOT NULL,
    "is_expert_solve" BOOLEAN NOT NULL,
    "fine_decision"   BOOLEAN NOT NULL,
    "got_at"          TIMESTAMP    NOT NULL,
    "solved_at"       TIMESTAMP
);
ALTER TABLE
    "expert_cases"
    ADD PRIMARY KEY ("expert_case_id");

ALTER TABLE
    "expert_cases"
    ADD CONSTRAINT "expert_cases_case_id_foreign" FOREIGN KEY ("case_id") REFERENCES "cases" ("case_id");
ALTER TABLE
    "expert_cases"
    ADD CONSTRAINT "expert_cases_expert_id_foreign" FOREIGN KEY ("expert_id") REFERENCES "experts" ("expert_id");