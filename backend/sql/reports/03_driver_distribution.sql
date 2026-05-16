SELECT
    e.id AS employee_id,
    e.last_name,
    e.first_name,
    v.id AS vehicle_id,
    v.license_plate,
    dva.start_date,
    dva.end_date
FROM driver_vehicle_assignments dva
JOIN drivers d ON d.employee_id = dva.driver_id
JOIN employees e ON e.id = d.employee_id
JOIN vehicles v ON v.id = dva.vehicle_id
ORDER BY v.license_plate, e.last_name, e.first_name;
