package application

import (
	"context"
	"fmt"

	"autopark/internal/domain"
)

type CRUDService struct {
	repo     domain.Repository
	entities map[string]domain.Entity
}

func NewCRUDService(repo domain.Repository) *CRUDService {
	return &CRUDService{repo: repo, entities: Entities()}
}

func (s *CRUDService) Entity(name string) (domain.Entity, error) {
	e, ok := s.entities[name]
	if !ok {
		return domain.Entity{}, fmt.Errorf("unknown entity %q", name)
	}
	return e, nil
}

func (s *CRUDService) Entities() map[string]domain.Entity {
	return s.entities
}

func (s *CRUDService) List(ctx context.Context, name string) ([]map[string]any, error) {
	e, err := s.Entity(name)
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, e)
}

func (s *CRUDService) Get(ctx context.Context, name string, id int64) (map[string]any, error) {
	e, err := s.Entity(name)
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, e, id)
}

func (s *CRUDService) Create(ctx context.Context, name string, values map[string]any) (map[string]any, error) {
	e, err := s.Entity(name)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, e, values)
}

func (s *CRUDService) Update(ctx context.Context, name string, id int64, values map[string]any) (map[string]any, error) {
	e, err := s.Entity(name)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, e, id, values)
}

func (s *CRUDService) Delete(ctx context.Context, name string, id int64) error {
	e, err := s.Entity(name)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, e, id)
}

type ReportService struct {
	repo domain.Repository
}

func NewReportService(repo domain.Repository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) Autopark(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ReportAutopark(ctx)
}

func (s *ReportService) Drivers(ctx context.Context, vehicleID any) ([]map[string]any, error) {
	return s.repo.ReportDrivers(ctx, vehicleID)
}

func (s *ReportService) DriverDistribution(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ReportDriverDistribution(ctx)
}

func (s *ReportService) PassengerRoutes(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ReportPassengerRoutes(ctx)
}

func (s *ReportService) Mileage(ctx context.Context, dateFrom, dateTo, categoryID, vehicleID any) ([]map[string]any, error) {
	return s.repo.ReportMileage(ctx, dateFrom, dateTo, categoryID, vehicleID)
}

func (s *ReportService) RepairsSummary(ctx context.Context, dateFrom, dateTo, categoryID, brand, vehicleID any) ([]map[string]any, error) {
	return s.repo.ReportRepairsSummary(ctx, dateFrom, dateTo, categoryID, brand, vehicleID)
}

func (s *ReportService) StaffHierarchy(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ReportStaffHierarchy(ctx)
}

func (s *ReportService) Garage(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ReportGarage(ctx)
}

func (s *ReportService) VehicleDistribution(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ReportVehicleDistribution(ctx)
}

func (s *ReportService) Freight(ctx context.Context, vehicleID, dateFrom, dateTo any) ([]map[string]any, error) {
	return s.repo.ReportFreight(ctx, vehicleID, dateFrom, dateTo)
}

func (s *ReportService) UsedParts(ctx context.Context, dateFrom, dateTo, categoryID, brand, vehicleID any) ([]map[string]any, error) {
	return s.repo.ReportUsedParts(ctx, dateFrom, dateTo, categoryID, brand, vehicleID)
}

func (s *ReportService) VehicleMovement(ctx context.Context, dateFrom, dateTo any) ([]map[string]any, error) {
	return s.repo.ReportVehicleMovement(ctx, dateFrom, dateTo)
}

func (s *ReportService) ManagerSubordinates(ctx context.Context, managerID any) ([]map[string]any, error) {
	return s.repo.ReportManagerSubordinates(ctx, managerID)
}

func (s *ReportService) RepairmanWorks(ctx context.Context, employeeID, dateFrom, dateTo, vehicleID any) ([]map[string]any, error) {
	return s.repo.ReportRepairmanWorks(ctx, employeeID, dateFrom, dateTo, vehicleID)
}
