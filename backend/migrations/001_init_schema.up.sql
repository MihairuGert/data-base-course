CREATE OR REPLACE TABLE positions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    level INTEGER NOT NULL CHECK (level > 0)
);

CREATE OR REPLACE TABLE garage_facilities (
    id BIGSERIAL PRIMARY KEY,
    facility_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(150) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity >= 0)
);

CREATE OR REPLACE TABLE vehicle_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE OR REPLACE TABLE brigades (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    foreman_id BIGINT UNIQUE,
    facility_id BIGINT NOT NULL REFERENCES garage_facilities(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE OR REPLACE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    last_name VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    birth_date DATE NOT NULL,
    hire_date DATE NOT NULL,
    position_id BIGINT NOT NULL REFERENCES positions(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    manager_id BIGINT REFERENCES employees(id) ON UPDATE CASCADE ON DELETE SET NULL,
    brigade_id BIGINT REFERENCES brigades(id) ON UPDATE CASCADE ON DELETE SET NULL
);

ALTER TABLE brigades
    ADD CONSTRAINT fk_brigades_foreman
    FOREIGN KEY (foreman_id) REFERENCES employees(id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE OR REPLACE TABLE drivers (
    employee_id BIGINT PRIMARY KEY REFERENCES employees(id) ON UPDATE CASCADE ON DELETE CASCADE,
    license_category VARCHAR(10) NOT NULL,
    driving_experience INTEGER NOT NULL CHECK (driving_experience >= 0)
);

CREATE OR REPLACE TABLE repairmen (
    employee_id BIGINT PRIMARY KEY REFERENCES employees(id) ON UPDATE CASCADE ON DELETE CASCADE,
    specialization VARCHAR(50) NOT NULL,
    rank INTEGER NOT NULL CHECK (rank > 0)
);

CREATE OR REPLACE TABLE vehicles (
    id BIGSERIAL PRIMARY KEY,
    license_plate VARCHAR(15) NOT NULL UNIQUE,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INTEGER NOT NULL CHECK (year >= 1900),
    acquisition_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL CHECK (status IN ('в эксплуатации', 'в ремонте', 'списан')),
    disposal_date DATE,
    category_id BIGINT NOT NULL REFERENCES vehicle_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    facility_id BIGINT REFERENCES garage_facilities(id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE OR REPLACE TABLE buses (
    vehicle_id BIGINT PRIMARY KEY REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    passenger_capacity INTEGER NOT NULL CHECK (passenger_capacity > 0)
);

CREATE OR REPLACE TABLE route_taxis (
    vehicle_id BIGINT PRIMARY KEY REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    passenger_capacity INTEGER NOT NULL CHECK (passenger_capacity > 0)
);

CREATE OR REPLACE TABLE taxis (
    vehicle_id BIGINT PRIMARY KEY REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    passenger_capacity INTEGER NOT NULL CHECK (passenger_capacity > 0)
);

CREATE OR REPLACE TABLE trucks (
    vehicle_id BIGINT PRIMARY KEY REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    load_capacity NUMERIC(10,2) NOT NULL CHECK (load_capacity > 0)
);

CREATE OR REPLACE TABLE aux_vehicles (
    vehicle_id BIGINT PRIMARY KEY REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    aux_type VARCHAR(50) NOT NULL
);

CREATE OR REPLACE TABLE routes (
    id BIGSERIAL PRIMARY KEY,
    route_number VARCHAR(20) NOT NULL UNIQUE,
    start_point VARCHAR(100) NOT NULL,
    end_point VARCHAR(100) NOT NULL,
    distance NUMERIC(10,2) NOT NULL CHECK (distance > 0)
);

CREATE OR REPLACE TABLE route_assignments (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    route_id BIGINT NOT NULL REFERENCES routes(id) ON UPDATE CASCADE ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE,
    note TEXT,
    CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE OR REPLACE TABLE driver_vehicle_assignments (
    id BIGSERIAL PRIMARY KEY,
    driver_id BIGINT NOT NULL REFERENCES drivers(employee_id) ON UPDATE CASCADE ON DELETE CASCADE,
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE,
    CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE OR REPLACE TABLE transport_logs (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    route_id BIGINT REFERENCES routes(id) ON UPDATE CASCADE ON DELETE SET NULL,
    log_date DATE NOT NULL,
    mileage NUMERIC(10,2) NOT NULL CHECK (mileage >= 0),
    passenger_count INTEGER CHECK (passenger_count >= 0),
    cargo_volume NUMERIC(10,2) CHECK (cargo_volume >= 0),
    note TEXT,
    UNIQUE (vehicle_id, log_date)
);

CREATE OR REPLACE TABLE repairs (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    brigade_id BIGINT REFERENCES brigades(id) ON UPDATE CASCADE ON DELETE SET NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    repair_type VARCHAR(30) NOT NULL CHECK (repair_type IN ('ТО', 'текущий', 'капитальный')),
    total_cost NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (total_cost >= 0),
    CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE OR REPLACE TABLE repair_works (
    id BIGSERIAL PRIMARY KEY,
    repair_id BIGINT NOT NULL REFERENCES repairs(id) ON UPDATE CASCADE ON DELETE CASCADE,
    employee_id BIGINT REFERENCES employees(id) ON UPDATE CASCADE ON DELETE SET NULL,
    work_type VARCHAR(100) NOT NULL,
    hours NUMERIC(8,2) NOT NULL CHECK (hours >= 0),
    cost NUMERIC(12,2) NOT NULL CHECK (cost >= 0)
);

CREATE OR REPLACE TABLE parts (
    id BIGSERIAL PRIMARY KEY,
    part_number VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    unit VARCHAR(20) NOT NULL
);

CREATE OR REPLACE TABLE replaced_parts (
    id BIGSERIAL PRIMARY KEY,
    repair_id BIGINT NOT NULL REFERENCES repairs(id) ON UPDATE CASCADE ON DELETE CASCADE,
    part_id BIGINT NOT NULL REFERENCES parts(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
    total_cost NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (total_cost >= 0)
);

CREATE OR REPLACE TABLE part_requests (
    id BIGSERIAL PRIMARY KEY,
    request_date DATE NOT NULL,
    brigade_id BIGINT REFERENCES brigades(id) ON UPDATE CASCADE ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('создана', 'в работе', 'выполнена', 'отменена'))
);

CREATE OR REPLACE TABLE part_request_items (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES part_requests(id) ON UPDATE CASCADE ON DELETE CASCADE,
    part_id BIGINT NOT NULL REFERENCES parts(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    note TEXT
);

CREATE INDEX idx_employees_manager_id ON employees(manager_id);
CREATE INDEX idx_vehicles_category_id ON vehicles(category_id);
CREATE INDEX idx_transport_logs_date ON transport_logs(log_date);
CREATE INDEX idx_repairs_vehicle_id ON repairs(vehicle_id);
CREATE INDEX idx_repairs_dates ON repairs(start_date, end_date);
