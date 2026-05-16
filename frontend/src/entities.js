const n = "number";
const d = "date";

export const entities = {
  positions: { title: "Должности", endpoint: "positions", fields: [
    ["name", "Название"], ["level", "Уровень", n],
  ] },
  facilities: { title: "Гаражное хозяйство", endpoint: "facilities", fields: [
    ["facility_type", "Тип"], ["name", "Название"], ["location", "Расположение"], ["capacity", "Вместимость", n],
  ] },
  "vehicle-categories": { title: "Категории транспорта", endpoint: "vehicle-categories", fields: [
    ["name", "Название"],
  ] },
  parts: { title: "Узлы и агрегаты", endpoint: "parts", fields: [
    ["part_number", "Артикул"], ["name", "Название"], ["category", "Категория"], ["unit", "Ед. изм."],
  ] },
  routes: { title: "Маршруты", endpoint: "routes", fields: [
    ["route_number", "Номер"], ["start_point", "Начало"], ["end_point", "Конец"], ["distance", "Дистанция", n],
  ] },
  employees: { title: "Сотрудники", endpoint: "employees", fields: [
    ["last_name", "Фамилия"], ["first_name", "Имя"], ["middle_name", "Отчество", "text", true],
    ["birth_date", "Дата рождения", d], ["hire_date", "Дата приема", d], ["position_id", "Должность ID", n],
    ["manager_id", "Руководитель ID", n, true], ["brigade_id", "Бригада ID", n, true],
  ] },
  drivers: { title: "Водители", endpoint: "drivers", id: "employee_id", fields: [
    ["employee_id", "Сотрудник ID", n], ["license_category", "Категория прав"], ["driving_experience", "Стаж", n],
  ] },
  repairmen: { title: "Ремонтники", endpoint: "repairmen", id: "employee_id", fields: [
    ["employee_id", "Сотрудник ID", n], ["specialization", "Специализация"], ["rank", "Разряд", n],
  ] },
  brigades: { title: "Бригады", endpoint: "brigades", fields: [
    ["name", "Название"], ["foreman_id", "Бригадир ID", n, true], ["facility_id", "Объект ID", n],
  ] },
  vehicles: { title: "Транспорт", endpoint: "vehicles", fields: [
    ["license_plate", "Госномер"], ["brand", "Марка"], ["model", "Модель"], ["year", "Год", n],
    ["acquisition_date", "Поступил", d], ["status", "Статус"], ["disposal_date", "Списан", d, true],
    ["category_id", "Категория ID", n], ["facility_id", "Объект ID", n, true],
  ] },
  buses: { title: "Автобусы", endpoint: "buses", id: "vehicle_id", fields: [
    ["vehicle_id", "ТС ID", n], ["passenger_capacity", "Пассажиров", n],
  ] },
  "route-taxis": { title: "Маршрутные такси", endpoint: "route-taxis", id: "vehicle_id", fields: [
    ["vehicle_id", "ТС ID", n], ["passenger_capacity", "Пассажиров", n],
  ] },
  taxis: { title: "Такси", endpoint: "taxis", id: "vehicle_id", fields: [
    ["vehicle_id", "ТС ID", n], ["passenger_capacity", "Пассажиров", n],
  ] },
  trucks: { title: "Грузовой транспорт", endpoint: "trucks", id: "vehicle_id", fields: [
    ["vehicle_id", "ТС ID", n], ["load_capacity", "Грузоподъемность", n],
  ] },
  "aux-vehicles": { title: "Вспомогательный транспорт", endpoint: "aux-vehicles", id: "vehicle_id", fields: [
    ["vehicle_id", "ТС ID", n], ["aux_type", "Тип"],
  ] },
  "route-assignments": { title: "Назначения на маршруты", endpoint: "route-assignments", fields: [
    ["vehicle_id", "ТС ID", n], ["route_id", "Маршрут ID", n], ["start_date", "Начало", d], ["end_date", "Конец", d, true], ["note", "Примечание", "text", true],
  ] },
  "driver-vehicle-assignments": { title: "Закрепление водителей", endpoint: "driver-vehicle-assignments", fields: [
    ["driver_id", "Водитель ID", n], ["vehicle_id", "ТС ID", n], ["start_date", "Начало", d], ["end_date", "Конец", d, true],
  ] },
  "transport-logs": { title: "Перевозки", endpoint: "transport-logs", fields: [
    ["vehicle_id", "ТС ID", n], ["route_id", "Маршрут ID", n, true], ["log_date", "Дата", d], ["mileage", "Пробег", n],
    ["passenger_count", "Пассажиры", n, true], ["cargo_volume", "Груз", n, true], ["note", "Примечание", "text", true],
  ] },
  repairs: { title: "Ремонты", endpoint: "repairs", fields: [
    ["vehicle_id", "ТС ID", n], ["brigade_id", "Бригада ID", n, true], ["start_date", "Начало", d], ["end_date", "Конец", d, true],
    ["repair_type", "Тип"], ["total_cost", "Стоимость", n, true],
  ] },
  "repair-works": { title: "Работы по ремонту", endpoint: "repair-works", fields: [
    ["repair_id", "Ремонт ID", n], ["employee_id", "Сотрудник ID", n, true], ["work_type", "Работа"], ["hours", "Часы", n], ["cost", "Стоимость", n],
  ] },
  "replaced-parts": { title: "Замененные детали", endpoint: "replaced-parts", fields: [
    ["repair_id", "Ремонт ID", n], ["part_id", "Деталь ID", n], ["quantity", "Количество", n], ["unit_price", "Цена", n], ["total_cost", "Итого", n, true],
  ] },
  "part-requests": { title: "Заявки", endpoint: "part-requests", fields: [
    ["request_date", "Дата", d], ["brigade_id", "Бригада ID", n, true], ["status", "Статус"],
  ] },
  "part-request-items": { title: "Позиции заявок", endpoint: "part-request-items", fields: [
    ["request_id", "Заявка ID", n], ["part_id", "Деталь ID", n], ["quantity", "Количество", n], ["note", "Примечание", "text", true],
  ] },
};

export const groups = {
  catalogs: ["positions", "facilities", "vehicle-categories", "parts", "routes"],
  employees: ["employees", "drivers", "repairmen", "brigades"],
  vehicles: ["vehicles", "buses", "route-taxis", "taxis", "trucks", "aux-vehicles", "driver-vehicle-assignments"],
  routes: ["routes", "route-assignments"],
  transport: ["transport-logs"],
  repairs: ["repairs", "repair-works", "replaced-parts"],
  requests: ["part-requests", "part-request-items"],
};
