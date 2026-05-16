WITH driver_rows AS (
    SELECT
        e.id AS employee_id,
        e.last_name,
        e.first_name,
        e.middle_name,
        d.license_category,
        d.driving_experience,
        v.id AS vehicle_id,
        v.license_plate
    FROM drivers d
    JOIN employees e ON e.id = d.employee_id
    LEFT JOIN driver_vehicle_assignments dva
        ON dva.driver_id = d.employee_id
       AND dva.end_date IS NULL
    LEFT JOIN vehicles v ON v.id = dva.vehicle_id
    WHERE ($1::bigint IS NULL OR v.id = $1)
)
SELECT *, COUNT(*) OVER () AS driver_count
FROM driver_rows
ORDER BY last_name, first_name;
