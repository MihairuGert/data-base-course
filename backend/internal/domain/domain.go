package domain

import "context"

type FieldKind string

const (
	Text    FieldKind = "text"
	Int     FieldKind = "int"
	Numeric FieldKind = "numeric"
	Date    FieldKind = "date"
)

type Field struct {
	Name     string
	Kind     FieldKind
	Nullable bool
}

type Entity struct {
	Name   string
	Table  string
	ID     string
	Order  string
	Fields []Field
}

func (e Entity) Columns() []string {
	seen := map[string]bool{e.ID: true}
	cols := []string{e.ID}
	for _, f := range e.Fields {
		if !seen[f.Name] {
			cols = append(cols, f.Name)
			seen[f.Name] = true
		}
	}
	return cols
}

type Repository interface {
	List(ctx context.Context, e Entity) ([]map[string]any, error)
	Get(ctx context.Context, e Entity, id int64) (map[string]any, error)
	Create(ctx context.Context, e Entity, values map[string]any) (map[string]any, error)
	Update(ctx context.Context, e Entity, id int64, values map[string]any) (map[string]any, error)
	Delete(ctx context.Context, e Entity, id int64) error

	ReportAutopark(ctx context.Context) ([]map[string]any, error)
	ReportDrivers(ctx context.Context, vehicleID any) ([]map[string]any, error)
	ReportDriverDistribution(ctx context.Context) ([]map[string]any, error)
	ReportPassengerRoutes(ctx context.Context) ([]map[string]any, error)
	ReportMileage(ctx context.Context, dateFrom, dateTo, categoryID, vehicleID any) ([]map[string]any, error)
	ReportRepairsSummary(ctx context.Context, dateFrom, dateTo, categoryID, brand, vehicleID any) ([]map[string]any, error)
	ReportStaffHierarchy(ctx context.Context) ([]map[string]any, error)
	ReportGarage(ctx context.Context) ([]map[string]any, error)
	ReportVehicleDistribution(ctx context.Context) ([]map[string]any, error)
	ReportFreight(ctx context.Context, vehicleID, dateFrom, dateTo any) ([]map[string]any, error)
	ReportUsedParts(ctx context.Context, dateFrom, dateTo, categoryID, brand, vehicleID any) ([]map[string]any, error)
	ReportVehicleMovement(ctx context.Context, dateFrom, dateTo any) ([]map[string]any, error)
	ReportManagerSubordinates(ctx context.Context, managerID any) ([]map[string]any, error)
	ReportRepairmanWorks(ctx context.Context, employeeID, dateFrom, dateTo, vehicleID any) ([]map[string]any, error)

	FindUserByUsername(ctx context.Context, username string) (User, error)
}

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         string
	EmployeeID   *int64
}
