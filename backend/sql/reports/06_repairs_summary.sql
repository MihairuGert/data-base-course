SELECT
    vc.name AS category_name,
    v.brand,
    v.id AS vehicle_id,
    v.license_plate,
    COUNT(r.id) AS repair_count,
    COALESCE(SUM(r.total_cost), 0) AS repair_total_cost
FROM repairs r
JOIN vehicles v ON v.id = r.vehicle_id
JOIN vehicle_categories vc ON vc.id = v.category_id
WHERE r.start_date BETWEEN COALESCE($1::date, DATE '1900-01-01') AND COALESCE($2::date, DATE '2999-12-31')
  AND ($3::bigint IS NULL OR v.category_id = $3)
  AND ($4::text IS NULL OR v.brand = $4)
  AND ($5::bigint IS NULL OR v.id = $5)
GROUP BY vc.name, v.brand, v.id, v.license_plate
ORDER BY vc.name, v.brand, v.license_plate;
