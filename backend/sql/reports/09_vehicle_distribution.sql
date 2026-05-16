SELECT
    v.id AS vehicle_id,
    v.license_plate,
    v.brand,
    v.model,
    vc.name AS category_name,
    gf.name AS facility_name,
    v.status
FROM vehicles v
JOIN vehicle_categories vc ON vc.id = v.category_id
LEFT JOIN garage_facilities gf ON gf.id = v.facility_id
ORDER BY gf.name, vc.name, v.license_plate;
