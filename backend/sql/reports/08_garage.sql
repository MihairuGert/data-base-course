SELECT
    gf.id AS facility_id,
    gf.name AS facility_name,
    gf.facility_type,
    gf.capacity,
    COALESCE(vc.name, 'все категории') AS category_name,
    COUNT(v.id) AS vehicle_count
FROM garage_facilities gf
LEFT JOIN vehicles v ON v.facility_id = gf.id
LEFT JOIN vehicle_categories vc ON vc.id = v.category_id
GROUP BY GROUPING SETS (
    (gf.id, gf.name, gf.facility_type, gf.capacity, vc.name),
    (gf.id, gf.name, gf.facility_type, gf.capacity)
)
ORDER BY gf.name, category_name;
