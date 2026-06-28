DROP TABLE todoapp.task_tags;
DROP TABLE todoapp.tags;

ALTER TABLE todoapp.tasks
    DROP COLUMN archived_at,
    DROP COLUMN due_at,
    DROP COLUMN priority;
