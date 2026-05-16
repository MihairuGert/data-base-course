WITH RECURSIVE emp_tree AS (
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
    WHERE e.manager_id IS NULL

    UNION ALL

    SELECT
        e.id,
        e.last_name,
        e.first_name,
        e.middle_name,
        p.name,
        e.manager_id,
        et.level + 1
    FROM employees e
    JOIN positions p ON p.id = e.position_id
    JOIN emp_tree et ON et.employee_id = e.manager_id
)
SELECT *
FROM emp_tree
ORDER BY level, last_name, first_name;
