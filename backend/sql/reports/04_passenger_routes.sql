SELECT
    v.id AS vehicle_id,
    v.license_plate,
    vc.name AS category_name,
    r.route_number,
    r.start_point,
    r.end_point,
    ra.start_date,
    ra.end_date
FROM route_assignments ra
JOIN vehicles v ON v.id = ra.vehicle_id
JOIN vehicle_categories vc ON vc.id = v.category_id
JOIN routes r ON r.id = ra.route_id
WHERE vc.name IN ('автобус', 'маршрутное такси', 'такси')
ORDER BY r.route_number, v.license_plate;
