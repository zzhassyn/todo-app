BEGIN;

DROP INDEX IF EXISTS todoapp.idx_tasks_position;

ALTER TABLE todoapp.tasks
DROP COLUMN position;

COMMIT;
