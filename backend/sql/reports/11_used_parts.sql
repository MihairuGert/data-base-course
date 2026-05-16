SELECT
    p.id AS part_id,
    p.name AS part_name,
    p.part_number,
    vc.name AS category_name,
    v.brand,
    v.id AS vehicle_id,
    v.license_plate,
    SUM(rp.quantity) AS total_used_quantity,
    SUM(rp.total_cost) AS total_cost
FROM replaced_parts rp
JOIN repairs r ON r.id = rp.repair_id
JOIN vehicles v ON v.id = r.vehicle_id
JOIN vehicle_categories vc ON vc.id = v.category_id
JOIN parts p ON p.id = rp.part_id
WHERE r.start_date BETWEEN COALESCE($1::date, DATE '1900-01-01') AND COALESCE($2::date, DATE '2999-12-31')
  AND ($3::bigint IS NULL OR v.category_id = $3)
  AND ($4::text IS NULL OR v.brand = $4)
  AND ($5::bigint IS NULL OR v.id = $5)
GROUP BY p.id, p.name, p.part_number, vc.name, v.brand, v.id, v.license_plate
ORDER BY p.name, v.license_plate;
