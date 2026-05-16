package application

import "autopark/internal/domain"

func Entities() map[string]domain.Entity {
	t := domain.Text
	i := domain.Int
	n := domain.Numeric
	d := domain.Date

	return map[string]domain.Entity{
		"positions": {Name: "positions", Table: "positions", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "name", Kind: t}, {Name: "level", Kind: i},
		}},
		"facilities": {Name: "facilities", Table: "garage_facilities", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "facility_type", Kind: t}, {Name: "name", Kind: t}, {Name: "location", Kind: t}, {Name: "capacity", Kind: i},
		}},
		"vehicle-categories": {Name: "vehicle-categories", Table: "vehicle_categories", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "name", Kind: t},
		}},
		"brigades": {Name: "brigades", Table: "brigades", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "name", Kind: t}, {Name: "foreman_id", Kind: i, Nullable: true}, {Name: "facility_id", Kind: i},
		}},
		"employees": {Name: "employees", Table: "employees", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "last_name", Kind: t}, {Name: "first_name", Kind: t}, {Name: "middle_name", Kind: t, Nullable: true},
			{Name: "birth_date", Kind: d}, {Name: "hire_date", Kind: d}, {Name: "position_id", Kind: i},
			{Name: "manager_id", Kind: i, Nullable: true}, {Name: "brigade_id", Kind: i, Nullable: true},
		}},
		"drivers": {Name: "drivers", Table: "drivers", ID: "employee_id", Order: "employee_id", Fields: []domain.Field{
			{Name: "employee_id", Kind: i}, {Name: "license_category", Kind: t}, {Name: "driving_experience", Kind: i},
		}},
		"repairmen": {Name: "repairmen", Table: "repairmen", ID: "employee_id", Order: "employee_id", Fields: []domain.Field{
			{Name: "employee_id", Kind: i}, {Name: "specialization", Kind: t}, {Name: "rank", Kind: i},
		}},
		"vehicles": {Name: "vehicles", Table: "vehicles", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "license_plate", Kind: t}, {Name: "brand", Kind: t}, {Name: "model", Kind: t}, {Name: "year", Kind: i},
			{Name: "acquisition_date", Kind: d}, {Name: "status", Kind: t}, {Name: "disposal_date", Kind: d, Nullable: true},
			{Name: "category_id", Kind: i}, {Name: "facility_id", Kind: i, Nullable: true},
		}},
		"buses": {Name: "buses", Table: "buses", ID: "vehicle_id", Order: "vehicle_id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "passenger_capacity", Kind: i},
		}},
		"route-taxis": {Name: "route-taxis", Table: "route_taxis", ID: "vehicle_id", Order: "vehicle_id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "passenger_capacity", Kind: i},
		}},
		"taxis": {Name: "taxis", Table: "taxis", ID: "vehicle_id", Order: "vehicle_id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "passenger_capacity", Kind: i},
		}},
		"trucks": {Name: "trucks", Table: "trucks", ID: "vehicle_id", Order: "vehicle_id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "load_capacity", Kind: n},
		}},
		"aux-vehicles": {Name: "aux-vehicles", Table: "aux_vehicles", ID: "vehicle_id", Order: "vehicle_id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "aux_type", Kind: t},
		}},
		"routes": {Name: "routes", Table: "routes", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "route_number", Kind: t}, {Name: "start_point", Kind: t}, {Name: "end_point", Kind: t}, {Name: "distance", Kind: n},
		}},
		"route-assignments": {Name: "route-assignments", Table: "route_assignments", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "route_id", Kind: i}, {Name: "start_date", Kind: d}, {Name: "end_date", Kind: d, Nullable: true}, {Name: "note", Kind: t, Nullable: true},
		}},
		"driver-vehicle-assignments": {Name: "driver-vehicle-assignments", Table: "driver_vehicle_assignments", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "driver_id", Kind: i}, {Name: "vehicle_id", Kind: i}, {Name: "start_date", Kind: d}, {Name: "end_date", Kind: d, Nullable: true},
		}},
		"transport-logs": {Name: "transport-logs", Table: "transport_logs", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "route_id", Kind: i, Nullable: true}, {Name: "log_date", Kind: d}, {Name: "mileage", Kind: n},
			{Name: "passenger_count", Kind: i, Nullable: true}, {Name: "cargo_volume", Kind: n, Nullable: true}, {Name: "note", Kind: t, Nullable: true},
		}},
		"repairs": {Name: "repairs", Table: "repairs", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "vehicle_id", Kind: i}, {Name: "brigade_id", Kind: i, Nullable: true}, {Name: "start_date", Kind: d}, {Name: "end_date", Kind: d, Nullable: true},
			{Name: "repair_type", Kind: t}, {Name: "total_cost", Kind: n, Nullable: true},
		}},
		"repair-works": {Name: "repair-works", Table: "repair_works", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "repair_id", Kind: i}, {Name: "employee_id", Kind: i, Nullable: true}, {Name: "work_type", Kind: t}, {Name: "hours", Kind: n}, {Name: "cost", Kind: n},
		}},
		"parts": {Name: "parts", Table: "parts", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "part_number", Kind: t}, {Name: "name", Kind: t}, {Name: "category", Kind: t}, {Name: "unit", Kind: t},
		}},
		"replaced-parts": {Name: "replaced-parts", Table: "replaced_parts", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "repair_id", Kind: i}, {Name: "part_id", Kind: i}, {Name: "quantity", Kind: i}, {Name: "unit_price", Kind: n}, {Name: "total_cost", Kind: n, Nullable: true},
		}},
		"part-requests": {Name: "part-requests", Table: "part_requests", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "request_date", Kind: d}, {Name: "brigade_id", Kind: i, Nullable: true}, {Name: "status", Kind: t},
		}},
		"part-request-items": {Name: "part-request-items", Table: "part_request_items", ID: "id", Order: "id", Fields: []domain.Field{
			{Name: "request_id", Kind: i}, {Name: "part_id", Kind: i}, {Name: "quantity", Kind: i}, {Name: "note", Kind: t, Nullable: true},
		}},
	}
}
