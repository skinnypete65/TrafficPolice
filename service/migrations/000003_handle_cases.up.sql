CREATE TABLE "solved_cases"
(
    "solved_case_id"  UUID    NOT NULL,
    "expert_id"       UUID    NOT NULL,
    "case_id"         UUID    NOT NULL,
    "is_expert_solve" BOOLEAN NOT NULL,
    "fine_decision"   BOOLEAN NOT NULL
);
ALTER TABLE
    "solved_cases"
    ADD PRIMARY KEY ("solved_case_id");

ALTER TABLE
    "solved_cases"
    ADD CONSTRAINT "solved_cases_case_id_foreign" FOREIGN KEY ("case_id") REFERENCES "cases" ("case_id");
ALTER TABLE
    "solved_cases"
    ADD CONSTRAINT "solved_cases_expert_id_foreign" FOREIGN KEY ("expert_id") REFERENCES "experts" ("expert_id");