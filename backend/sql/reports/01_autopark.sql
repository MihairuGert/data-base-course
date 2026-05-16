SELECT
    v.id AS vehicle_id,
    v.license_plate,
    v.brand,
    v.model,
    vc.name AS category_name,
    v.year,
    v.acquisition_date,
    v.status,
    gf.name AS facility_name
FROM vehicles v
JOIN vehicle_categories vc ON vc.id = v.category_id
LEFT JOIN garage_facilities gf ON gf.id = v.facility_id
ORDER BY vc.name, v.brand, v.model, v.license_plate;
