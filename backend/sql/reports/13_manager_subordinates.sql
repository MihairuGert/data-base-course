WITH RECURSIVE subordinates AS (
    SELECT
        e.id AS employee_id,
        e.last_name,
        e.first_name,
        e.middle_name,
        p.name AS position_name,
        e.manager_id,
        1 AS level
    FROM employees e
    JOIN positions p ON p.id = e.position_id
    WHERE e.manager_id = $1::bigint

    UNION ALL

    SELECT
        e.id,
        e.last_name,
        e.first_name,
        e.middle_name,
        p.name,
        e.manager_id,
        s.level + 1
    FROM employees e
    JOIN positions p ON p.id = e.position_id
    JOIN subordinates s ON s.employee_id = e.manager_id
)
SELECT *
FROM subordinates
ORDER BY level, last_name, first_name;
