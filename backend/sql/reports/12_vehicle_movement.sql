SELECT
    id AS vehicle_id,
    license_plate,
    brand,
    model,
    acquisition_date,
    disposal_date,
    status
FROM vehicles
WHERE acquisition_date BETWEEN COALESCE($1::date, DATE '1900-01-01') AND COALESCE($2::date, DATE '2999-12-31')
   OR disposal_date BETWEEN COALESCE($1::date, DATE '1900-01-01') AND COALESCE($2::date, DATE '2999-12-31')
ORDER BY acquisition_date, disposal_date;
