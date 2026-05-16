SELECT
    e.id AS employee_id,
    e.last_name,
    e.first_name,
    rw.id AS work_id,
    rw.work_type,
    rw.hours,
    rw.cost,
    r.id AS repair_id,
    v.id AS vehicle_id,
    v.license_plate,
    r.start_date,
    r.end_date
FROM repair_works rw
JOIN employees e ON e.id = rw.employee_id
JOIN repairs r ON r.id = rw.repair_id
JOIN vehicles v ON v.id = r.vehicle_id
WHERE ($1::bigint IS NULL OR rw.employee_id = $1)
  AND r.start_date BETWEEN COALESCE($2::date, DATE '1900-01-01') AND COALESCE($3::date, DATE '2999-12-31')
  AND ($4::bigint IS NULL OR v.id = $4)
ORDER BY r.start_date, rw.id;
