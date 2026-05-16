package postgres

import (
	"context"

	sqlfiles "autopark/sql"
)

func (s *Store) report(ctx context.Context, file string, args ...any) ([]map[string]any, error) {
	return s.query(ctx, sqlfiles.MustRead(file), args...)
}

func (s *Store) ReportAutopark(ctx context.Context) ([]map[string]any, error) {
	return s.report(ctx, "reports/01_autopark.sql")
}

func (s *Store) ReportDrivers(ctx context.Context, vehicleID any) ([]map[string]any, error) {
	return s.report(ctx, "reports/02_drivers.sql", vehicleID)
}

func (s *Store) ReportDriverDistribution(ctx context.Context) ([]map[string]any, error) {
	return s.report(ctx, "reports/03_driver_distribution.sql")
}

func (s *Store) ReportPassengerRoutes(ctx context.Context) ([]map[string]any, error) {
	return s.report(ctx, "reports/04_passenger_routes.sql")
}

func (s *Store) ReportMileage(ctx context.Context, dateFrom, dateTo, categoryID, vehicleID any) ([]map[string]any, error) {
	return s.report(ctx, "reports/05_mileage.sql", dateFrom, dateTo, categoryID, vehicleID)
}

func (s *Store) ReportRepairsSummary(ctx context.Context, dateFrom, dateTo, categoryID, brand, vehicleID any) ([]map[string]any, error) {
	return s.report(ctx, "reports/06_repairs_summary.sql", dateFrom, dateTo, categoryID, brand, vehicleID)
}

func (s *Store) ReportStaffHierarchy(ctx context.Context) ([]map[string]any, error) {
	return s.report(ctx, "reports/07_staff_hierarchy.sql")
}

func (s *Store) ReportGarage(ctx context.Context) ([]map[string]any, error) {
	return s.report(ctx, "reports/08_garage.sql")
}

func (s *Store) ReportVehicleDistribution(ctx context.Context) ([]map[string]any, error) {
	return s.report(ctx, "reports/09_vehicle_distribution.sql")
}

func (s *Store) ReportFreight(ctx context.Context, vehicleID, dateFrom, dateTo any) ([]map[string]any, error) {
	return s.report(ctx, "reports/10_freight.sql", vehicleID, dateFrom, dateTo)
}

func (s *Store) ReportUsedParts(ctx context.Context, dateFrom, dateTo, categoryID, brand, vehicleID any) ([]map[string]any, error) {
	return s.report(ctx, "reports/11_used_parts.sql", dateFrom, dateTo, categoryID, brand, vehicleID)
}

func (s *Store) ReportVehicleMovement(ctx context.Context, dateFrom, dateTo any) ([]map[string]any, error) {
	return s.report(ctx, "reports/12_vehicle_movement.sql", dateFrom, dateTo)
}

func (s *Store) ReportManagerSubordinates(ctx context.Context, managerID any) ([]map[string]any, error) {
	return s.report(ctx, "reports/13_manager_subordinates.sql", managerID)
}

func (s *Store) ReportRepairmanWorks(ctx context.Context, employeeID, dateFrom, dateTo, vehicleID any) ([]map[string]any, error) {
	return s.report(ctx, "reports/14_repairman_works.sql", employeeID, dateFrom, dateTo, vehicleID)
}
