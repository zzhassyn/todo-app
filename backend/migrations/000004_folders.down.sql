DROP INDEX todoapp.tasks_folder_id_idx;

ALTER TABLE todoapp.tasks
    DROP COLUMN folder_id;

DROP INDEX todoapp.folders_user_id_idx;
DROP TABLE todoapp.folders;
