ALTER TABLE student_activities
    ADD COLUMN is_embedded BOOLEAN DEFAULT FALSE;

CREATE INDEX idx_activities_is_embedded
    ON student_activities(is_embedded)
    WHERE is_embedded = FALSE;