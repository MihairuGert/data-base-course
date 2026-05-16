SELECT
    vc.name AS category_name,
    v.id AS vehicle_id,
    v.license_plate,
    SUM(tl.mileage) AS total_mileage
FROM transport_logs tl
JOIN vehicles v ON v.id = tl.vehicle_id
JOIN vehicle_categories vc ON vc.id = v.category_id
WHERE tl.log_date BETWEEN COALESCE($1::date, DATE '1900-01-01') AND COALESCE($2::date, DATE '2999-12-31')
  AND ($3::bigint IS NULL OR v.category_id = $3)
  AND ($4::bigint IS NULL OR v.id = $4)
GROUP BY vc.name, v.id, v.license_plate
ORDER BY vc.name, v.license_plate;
