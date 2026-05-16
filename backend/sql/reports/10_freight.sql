SELECT
    v.id AS vehicle_id,
    v.license_plate,
    tl.log_date,
    tl.mileage,
    tl.cargo_volume,
    tl.note
FROM transport_logs tl
JOIN vehicles v ON v.id = tl.vehicle_id
WHERE ($1::bigint IS NULL OR v.id = $1)
  AND tl.log_date BETWEEN COALESCE($2::date, DATE '1900-01-01') AND COALESCE($3::date, DATE '2999-12-31')
  AND tl.cargo_volume IS NOT NULL
ORDER BY tl.log_date;
