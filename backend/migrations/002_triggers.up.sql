CREATE OR REPLACE FUNCTION fn_set_replaced_part_total_cost()
RETURNS TRIGGER AS $$
BEGIN
    NEW.total_cost := NEW.quantity * NEW.unit_price;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_set_replaced_part_total_cost
BEFORE INSERT OR UPDATE ON replaced_parts
FOR EACH ROW EXECUTE FUNCTION fn_set_replaced_part_total_cost();

CREATE OR REPLACE FUNCTION fn_refresh_repair_total_cost(p_repair_id BIGINT)
RETURNS VOID AS $$
BEGIN
    UPDATE repairs r
    SET total_cost =
        COALESCE((SELECT SUM(cost) FROM repair_works WHERE repair_id = p_repair_id), 0) +
        COALESCE((SELECT SUM(total_cost) FROM replaced_parts WHERE repair_id = p_repair_id), 0)
    WHERE r.id = p_repair_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_repair_total_cost_trigger()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP IN ('UPDATE', 'DELETE') THEN
        PERFORM fn_refresh_repair_total_cost(OLD.repair_id);
    END IF;
    IF TG_OP IN ('INSERT', 'UPDATE') THEN
        PERFORM fn_refresh_repair_total_cost(NEW.repair_id);
    END IF;
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_refresh_repair_total_cost_from_work
AFTER INSERT OR UPDATE OR DELETE ON repair_works
FOR EACH ROW EXECUTE FUNCTION fn_repair_total_cost_trigger();

CREATE OR REPLACE TRIGGER trg_refresh_repair_total_cost_from_part
AFTER INSERT OR UPDATE OR DELETE ON replaced_parts
FOR EACH ROW EXECUTE FUNCTION fn_repair_total_cost_trigger();

CREATE OR REPLACE FUNCTION fn_check_route_assignment_overlap()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM route_assignments ra
        WHERE ra.vehicle_id = NEW.vehicle_id
          AND ra.id <> COALESCE(NEW.id, -1)
          AND daterange(ra.start_date, COALESCE(ra.end_date, 'infinity'::date), '[]')
              && daterange(NEW.start_date, COALESCE(NEW.end_date, 'infinity'::date), '[]')
    ) THEN
        RAISE EXCEPTION 'Для данного ТС уже существует пересекающееся назначение на маршрут';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_check_route_assignment_overlap
BEFORE INSERT OR UPDATE ON route_assignments
FOR EACH ROW EXECUTE FUNCTION fn_check_route_assignment_overlap();

CREATE OR REPLACE FUNCTION fn_check_driver_vehicle_overlap()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM driver_vehicle_assignments dva
        WHERE dva.driver_id = NEW.driver_id
          AND dva.id <> COALESCE(NEW.id, -1)
          AND daterange(dva.start_date, COALESCE(dva.end_date, 'infinity'::date), '[]')
              && daterange(NEW.start_date, COALESCE(NEW.end_date, 'infinity'::date), '[]')
    ) THEN
        RAISE EXCEPTION 'Для данного водителя уже существует пересекающееся закрепление';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_check_driver_vehicle_overlap
BEFORE INSERT OR UPDATE ON driver_vehicle_assignments
FOR EACH ROW EXECUTE FUNCTION fn_check_driver_vehicle_overlap();

CREATE OR REPLACE FUNCTION fn_validate_transport_log()
RETURNS TRIGGER AS $$
DECLARE
    v_category_name VARCHAR(50);
BEGIN
    SELECT vc.name INTO v_category_name
    FROM vehicles v
    JOIN vehicle_categories vc ON vc.id = v.category_id
    WHERE v.id = NEW.vehicle_id;

    IF v_category_name IN ('автобус', 'маршрутное такси', 'такси')
       AND COALESCE(NEW.cargo_volume, 0) > 0 THEN
        RAISE EXCEPTION 'Для пассажирского транспорта нельзя указывать объем груза';
    END IF;

    IF v_category_name = 'грузовой'
       AND COALESCE(NEW.passenger_count, 0) > 0 THEN
        RAISE EXCEPTION 'Для грузового транспорта нельзя указывать пассажиров как основной показатель';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_validate_transport_log
BEFORE INSERT OR UPDATE ON transport_logs
FOR EACH ROW EXECUTE FUNCTION fn_validate_transport_log();

CREATE OR REPLACE FUNCTION fn_repair_vehicle_status()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.end_date IS NULL THEN
            UPDATE vehicles SET status = 'в ремонте' WHERE id = NEW.vehicle_id AND status <> 'списан';
        END IF;
        RETURN NEW;
    END IF;

    IF TG_OP = 'UPDATE' THEN
        IF OLD.end_date IS NULL AND NEW.end_date IS NOT NULL THEN
            IF NOT EXISTS (
                SELECT 1 FROM repairs r
                WHERE r.vehicle_id = NEW.vehicle_id
                  AND r.id <> NEW.id
                  AND r.end_date IS NULL
            ) THEN
                UPDATE vehicles SET status = 'в эксплуатации' WHERE id = NEW.vehicle_id AND status <> 'списан';
            END IF;
        END IF;
        RETURN NEW;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_repair_vehicle_status_ins
AFTER INSERT ON repairs
FOR EACH ROW EXECUTE FUNCTION fn_repair_vehicle_status();

CREATE OR REPLACE TRIGGER trg_repair_vehicle_status_upd
AFTER UPDATE OF end_date ON repairs
FOR EACH ROW EXECUTE FUNCTION fn_repair_vehicle_status();

CREATE OR REPLACE FUNCTION fn_check_employee_manager_cycle()
RETURNS TRIGGER AS $$
DECLARE
    v_current BIGINT;
BEGIN
    v_current := NEW.manager_id;
    WHILE v_current IS NOT NULL LOOP
        IF v_current = NEW.id THEN
            RAISE EXCEPTION 'Обнаружена циклическая подчиненность';
        END IF;
        SELECT manager_id INTO v_current FROM employees WHERE id = v_current;
    END LOOP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_check_employee_manager_cycle
BEFORE INSERT OR UPDATE OF manager_id ON employees
FOR EACH ROW EXECUTE FUNCTION fn_check_employee_manager_cycle();

CREATE OR REPLACE FUNCTION fn_check_part_request_status()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status = 'выполнена' AND NEW.status <> 'выполнена' THEN
        RAISE EXCEPTION 'Нельзя менять статус уже выполненной заявки';
    END IF;
    IF OLD.status = 'отменена' AND NEW.status <> 'отменена' THEN
        RAISE EXCEPTION 'Нельзя менять статус отмененной заявки';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_check_part_request_status
BEFORE UPDATE OF status ON part_requests
FOR EACH ROW EXECUTE FUNCTION fn_check_part_request_status();
