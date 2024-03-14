CREATE TABLE "rating"
(
    "expert_id"     UUID   NOT NULL,
    "correct_cnt"   BIGINT NOT NULL,
    "incorrect_cnt" BIGINT NOT NULL
);
ALTER TABLE
    "rating"
    ADD PRIMARY KEY ("expert_id");

ALTER TABLE
    "rating"
    ADD CONSTRAINT "rating_expert_id_foreign" FOREIGN KEY ("expert_id") REFERENCES "experts" ("expert_id");