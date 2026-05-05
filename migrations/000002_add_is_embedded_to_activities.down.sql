DROP INDEX IF EXISTS idx_activities_is_embedded;
ALTER TABLE student_activities DROP COLUMN IF EXISTS is_embedded;