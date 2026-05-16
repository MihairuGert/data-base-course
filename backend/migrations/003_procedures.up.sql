CREATE OR REPLACE PROCEDURE sp_add_vehicle(
    p_license_plate VARCHAR(15),
    p_brand VARCHAR(50),
    p_model VARCHAR(50),
    p_year INTEGER,
    p_acquisition_date DATE,
    p_status VARCHAR(30),
    p_category_id BIGINT,
    p_facility_id BIGINT
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO vehicles (
        license_plate, brand, model, year, acquisition_date, status, category_id, facility_id
    )
    VALUES (
        p_license_plate, p_brand, p_model, p_year, p_acquisition_date, p_status, p_category_id, p_facility_id
    );
END;
$$;

CREATE OR REPLACE PROCEDURE sp_assign_vehicle_to_route(
    p_vehicle_id BIGINT,
    p_route_id BIGINT,
    p_start_date DATE,
    p_end_date DATE,
    p_note TEXT
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO route_assignments (vehicle_id, route_id, start_date, end_date, note)
    VALUES (p_vehicle_id, p_route_id, p_start_date, p_end_date, p_note);
END;
$$;

CREATE OR REPLACE PROCEDURE sp_add_transport_log(
    p_vehicle_id BIGINT,
    p_route_id BIGINT,
    p_log_date DATE,
    p_mileage NUMERIC(10,2),
    p_passenger_count INTEGER,
    p_cargo_volume NUMERIC(10,2),
    p_note TEXT
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO transport_logs (
        vehicle_id, route_id, log_date, mileage, passenger_count, cargo_volume, note
    )
    VALUES (
        p_vehicle_id, p_route_id, p_log_date, p_mileage, p_passenger_count, p_cargo_volume, p_note
    );
END;
$$;

CREATE OR REPLACE PROCEDURE sp_open_repair(
    p_vehicle_id BIGINT,
    p_brigade_id BIGINT,
    p_start_date DATE,
    p_repair_type VARCHAR(30)
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO repairs (vehicle_id, brigade_id, start_date, repair_type)
    VALUES (p_vehicle_id, p_brigade_id, p_start_date, p_repair_type);
END;
$$;

CREATE OR REPLACE PROCEDURE sp_close_repair(
    p_repair_id BIGINT,
    p_end_date DATE
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE repairs SET end_date = p_end_date WHERE id = p_repair_id;
    PERFORM fn_refresh_repair_total_cost(p_repair_id);
END;
$$;

CREATE OR REPLACE PROCEDURE sp_add_repair_work(
    p_repair_id BIGINT,
    p_employee_id BIGINT,
    p_work_type VARCHAR(100),
    p_hours NUMERIC(8,2),
    p_cost NUMERIC(12,2)
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO repair_works (repair_id, employee_id, work_type, hours, cost)
    VALUES (p_repair_id, p_employee_id, p_work_type, p_hours, p_cost);
END;
$$;

CREATE OR REPLACE PROCEDURE sp_add_replaced_part(
    p_repair_id BIGINT,
    p_part_id BIGINT,
    p_quantity INTEGER,
    p_unit_price NUMERIC(12,2)
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO replaced_parts (repair_id, part_id, quantity, unit_price)
    VALUES (p_repair_id, p_part_id, p_quantity, p_unit_price);
END;
$$;

CREATE OR REPLACE PROCEDURE sp_create_part_request(
    p_request_date DATE,
    p_brigade_id BIGINT,
    p_status VARCHAR(20)
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO part_requests (request_date, brigade_id, status)
    VALUES (p_request_date, p_brigade_id, p_status);
END;
$$;

CREATE OR REPLACE PROCEDURE sp_change_part_request_status(
    p_request_id BIGINT,
    p_status VARCHAR(20)
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE part_requests SET status = p_status WHERE id = p_request_id;
END;
$$;

CREATE OR REPLACE FUNCTION fn_get_vehicle_mileage(
    p_vehicle_id BIGINT,
    p_date_from DATE,
    p_date_to DATE
)
RETURNS NUMERIC(12,2)
LANGUAGE plpgsql AS $$
DECLARE
    v_total NUMERIC(12,2);
BEGIN
    SELECT COALESCE(SUM(mileage), 0)
    INTO v_total
    FROM transport_logs
    WHERE vehicle_id = p_vehicle_id
      AND log_date BETWEEN p_date_from AND p_date_to;

    RETURN v_total;
END;
$$;

CREATE OR REPLACE FUNCTION fn_report_autopark()
RETURNS TABLE (
    vehicle_id BIGINT,
    license_plate VARCHAR,
    brand VARCHAR,
    model VARCHAR,
    category_name VARCHAR,
    year INTEGER,
    acquisition_date DATE,
    status VARCHAR,
    facility_name VARCHAR
)
LANGUAGE sql AS $$
    SELECT
        v.id,
        v.license_plate,
        v.brand,
        v.model,
        vc.name,
        v.year,
        v.acquisition_date,
        v.status,
        gf.name
    FROM vehicles v
    JOIN vehicle_categories vc ON vc.id = v.category_id
    LEFT JOIN garage_facilities gf ON gf.id = v.facility_id
    ORDER BY vc.name, v.brand, v.model, v.license_plate;
$$;
