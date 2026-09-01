BEGIN;

DROP INDEX IF EXISTS todoapp.idx_tasks_search_vector;
ALTER TABLE todoapp.tasks DROP COLUMN IF EXISTS search_vector;

COMMIT;
