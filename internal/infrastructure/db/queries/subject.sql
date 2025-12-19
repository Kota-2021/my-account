-- name: SaveSubject :exec
INSERT INTO m_subjects (subject_code, subject_name)
VALUES ($1, $2)
ON CONFLICT (subject_code) DO UPDATE 
SET subject_name = EXCLUDED.subject_name;

-- name: ListSubjects :many
SELECT subject_code, subject_name 
FROM m_subjects 
ORDER BY subject_code;