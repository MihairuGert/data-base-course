INSERT INTO positions (name, level) VALUES
('водитель', 1),
('ремонтный рабочий', 1),
('бригадир', 2),
('мастер', 3),
('начальник участка', 4),
('директор', 5);

INSERT INTO garage_facilities (facility_type, name, location, capacity) VALUES
('гараж', 'Гараж №1', 'ул. Центральная, 10', 40),
('цех', 'Ремонтный цех №1', 'ул. Промышленная, 5', 20),
('бокс', 'Бокс диагностики', 'ул. Станционная, 7', 8);

INSERT INTO vehicle_categories (name) VALUES
('автобус'),
('маршрутное такси'),
('такси'),
('грузовой'),
('вспомогательный');

INSERT INTO brigades (name, facility_id) VALUES
('Бригада №1', 2),
('Бригада №2', 2),
('Водительская бригада', 1);

INSERT INTO employees (
    last_name, first_name, middle_name, birth_date, hire_date, position_id, manager_id, brigade_id
) VALUES
('Смирнов', 'Алексей', 'Игоревич', '1982-11-05', '2019-09-01', 6, NULL, NULL),
('Иванов', 'Иван', 'Иванович', '1985-02-10', '2020-03-01', 3, 1, 1),
('Петров', 'Петр', 'Петрович', '1990-04-12', '2021-06-15', 1, 2, 3),
('Сидоров', 'Сидор', 'Сидорович', '1988-08-20', '2022-01-10', 2, 2, 1),
('Кузнецова', 'Анна', 'Олеговна', '1992-07-18', '2022-03-14', 1, 2, 3),
('Орлов', 'Максим', 'Викторович', '1987-12-01', '2020-10-02', 2, 2, 2),
('Морозов', 'Денис', 'Сергеевич', '1995-09-09', '2023-05-12', 1, 2, 3);

UPDATE brigades SET foreman_id = 2 WHERE id = 1;
UPDATE brigades SET foreman_id = 6 WHERE id = 2;

INSERT INTO drivers (employee_id, license_category, driving_experience) VALUES
(3, 'D', 8),
(5, 'B', 5),
(7, 'C', 6);

INSERT INTO repairmen (employee_id, specialization, rank) VALUES
(4, 'слесарь', 4),
(6, 'электрик', 3);

INSERT INTO vehicles (
    license_plate, brand, model, year, acquisition_date, status, disposal_date, category_id, facility_id
) VALUES
('А123АА154', 'ПАЗ', '3205', 2018, '2023-01-15', 'в эксплуатации', NULL, 1, 1),
('М222ММ154', 'ГАЗ', 'ГАЗель Next', 2021, '2024-03-01', 'в эксплуатации', NULL, 2, 1),
('Т333ТТ154', 'Hyundai', 'Solaris', 2020, '2024-04-10', 'в эксплуатации', NULL, 3, 1),
('В456ВВ154', 'КАМАЗ', '65115', 2020, '2024-02-10', 'в эксплуатации', NULL, 4, 1),
('С789СС154', 'ГАЗ', '3309', 2017, '2022-05-20', 'в эксплуатации', NULL, 5, 3),
('О111ОО154', 'ЛиАЗ', '5292', 2012, '2018-06-11', 'списан', '2026-02-01', 1, NULL);

INSERT INTO buses (vehicle_id, passenger_capacity) VALUES
(1, 42),
(6, 95);

INSERT INTO route_taxis (vehicle_id, passenger_capacity) VALUES
(2, 18);

INSERT INTO taxis (vehicle_id, passenger_capacity) VALUES
(3, 4);

INSERT INTO trucks (vehicle_id, load_capacity) VALUES
(4, 12.50);

INSERT INTO aux_vehicles (vehicle_id, aux_type) VALUES
(5, 'эвакуатор');

INSERT INTO routes (route_number, start_point, end_point, distance) VALUES
('12', 'Площадь Ленина', 'ЖК Родники', 18.50),
('25', 'Вокзал', 'Академгородок', 24.20),
('7К', 'Речной вокзал', 'Затулинка', 15.10);

INSERT INTO route_assignments (vehicle_id, route_id, start_date, end_date, note) VALUES
(1, 1, '2026-01-01', NULL, 'Основной маршрут'),
(2, 2, '2026-01-10', NULL, 'Утренние рейсы'),
(3, 3, '2026-02-01', '2026-04-15', 'Пилотный маршрут'),
(6, 1, '2025-01-01', '2025-12-31', 'До списания');

INSERT INTO driver_vehicle_assignments (driver_id, vehicle_id, start_date, end_date) VALUES
(3, 1, '2026-01-01', NULL),
(5, 3, '2026-02-01', NULL),
(7, 4, '2026-01-15', NULL);

INSERT INTO transport_logs (
    vehicle_id, route_id, log_date, mileage, passenger_count, cargo_volume, note
) VALUES
(1, 1, '2026-04-01', 180.50, 950, NULL, 'Дневная смена'),
(1, 1, '2026-04-02', 175.20, 910, NULL, 'Дневная смена'),
(2, 2, '2026-04-01', 140.00, 360, NULL, 'Рейсы по расписанию'),
(3, 3, '2026-04-01', 95.00, 31, NULL, 'Таксомоторная смена'),
(4, NULL, '2026-04-01', 120.00, NULL, 18.30, 'Доставка стройматериалов'),
(4, NULL, '2026-04-02', 135.00, NULL, 21.10, 'Доставка материалов'),
(5, NULL, '2026-04-01', 40.00, NULL, NULL, 'Эвакуация транспорта');

INSERT INTO repairs (vehicle_id, brigade_id, start_date, end_date, repair_type) VALUES
(1, 1, '2026-03-10', '2026-03-12', 'текущий'),
(4, 2, '2026-04-10', NULL, 'капитальный');

INSERT INTO parts (part_number, name, category, unit) VALUES
('BRK-001', 'Тормозные колодки', 'тормозная система', 'комплект'),
('OIL-010', 'Масло моторное', 'двигатель', 'литр'),
('EL-020', 'Генератор', 'электрика', 'штука'),
('GBX-100', 'КПП', 'трансмиссия', 'штука');

INSERT INTO repair_works (repair_id, employee_id, work_type, hours, cost) VALUES
(1, 4, 'Замена тормозных колодок', 3.50, 4500.00),
(1, 6, 'Проверка электрики', 1.20, 1800.00),
(2, 4, 'Диагностика коробки передач', 4.00, 6000.00);

INSERT INTO replaced_parts (repair_id, part_id, quantity, unit_price) VALUES
(1, 1, 1, 3800.00),
(1, 2, 5, 650.00),
(2, 4, 1, 45000.00);

INSERT INTO part_requests (request_date, brigade_id, status) VALUES
('2026-04-10', 1, 'создана'),
('2026-04-11', 2, 'в работе');

INSERT INTO part_request_items (request_id, part_id, quantity, note) VALUES
(1, 1, 2, 'Для текущего ремонта'),
(1, 2, 10, 'Для планового ТО'),
(2, 3, 1, 'Для ремонта электрики'),
(2, 4, 1, 'Для капитального ремонта');
