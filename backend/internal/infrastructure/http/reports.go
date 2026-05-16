package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type reportInfo struct {
	Slug   string   `json:"slug"`
	Title  string   `json:"title"`
	Params []string `json:"params"`
}

var reports = []reportInfo{
	{"autopark", "Данные об автопарке предприятия", nil},
	{"drivers", "Перечень и общее число водителей", []string{"vehicle_id"}},
	{"driver-distribution", "Распределение водителей по автомобилям", nil},
	{"passenger-routes", "Распределение пассажирского транспорта по маршрутам", nil},
	{"mileage", "Пробег транспорта за период", []string{"date_from", "date_to", "category_id", "vehicle_id"}},
	{"repairs-summary", "Число ремонтов и их стоимость", []string{"date_from", "date_to", "category_id", "brand", "vehicle_id"}},
	{"staff-hierarchy", "Подчиненность персонала", nil},
	{"garage", "Наличие гаражного хозяйства", nil},
	{"vehicle-distribution", "Распределение автотранспорта", nil},
	{"freight", "Грузоперевозки машиной за период", []string{"vehicle_id", "date_from", "date_to"}},
	{"used-parts", "Использованные для ремонта узлы и агрегаты", []string{"date_from", "date_to", "category_id", "brand", "vehicle_id"}},
	{"vehicle-movement", "Полученная и списанная автотехника", []string{"date_from", "date_to"}},
	{"manager-subordinates", "Состав подчиненных руководителя", []string{"manager_id"}},
	{"repairman-works", "Работы специалиста за период", []string{"employee_id", "date_from", "date_to", "vehicle_id"}},
}

func (h *Handler) reportList(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	visible := make([]reportInfo, 0, len(reports))
	for _, report := range reports {
		if user.Permissions.Reports[report.Slug] {
			visible = append(visible, report)
		}
	}
	writeJSON(w, http.StatusOK, visible)
}

func (h *Handler) reportAutopark(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "autopark") {
		return
	}
	rows, err := h.reports.Autopark(r.Context())
	h.writeReport(w, rows, err)
}

func (h *Handler) reportDrivers(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "drivers") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.Drivers(r.Context(), optInt(q.Get("vehicle_id")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportDriverDistribution(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "driver-distribution") {
		return
	}
	rows, err := h.reports.DriverDistribution(r.Context())
	h.writeReport(w, rows, err)
}

func (h *Handler) reportPassengerRoutes(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "passenger-routes") {
		return
	}
	rows, err := h.reports.PassengerRoutes(r.Context())
	h.writeReport(w, rows, err)
}

func (h *Handler) reportMileage(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "mileage") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.Mileage(r.Context(), optDate(q.Get("date_from")), optDate(q.Get("date_to")), optInt(q.Get("category_id")), optInt(q.Get("vehicle_id")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportRepairsSummary(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "repairs-summary") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.RepairsSummary(r.Context(), optDate(q.Get("date_from")), optDate(q.Get("date_to")), optInt(q.Get("category_id")), optText(q.Get("brand")), optInt(q.Get("vehicle_id")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportStaffHierarchy(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "staff-hierarchy") {
		return
	}
	rows, err := h.reports.StaffHierarchy(r.Context())
	h.writeReport(w, rows, err)
}

func (h *Handler) reportGarage(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "garage") {
		return
	}
	rows, err := h.reports.Garage(r.Context())
	h.writeReport(w, rows, err)
}

func (h *Handler) reportVehicleDistribution(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "vehicle-distribution") {
		return
	}
	rows, err := h.reports.VehicleDistribution(r.Context())
	h.writeReport(w, rows, err)
}

func (h *Handler) reportFreight(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "freight") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.Freight(r.Context(), optInt(q.Get("vehicle_id")), optDate(q.Get("date_from")), optDate(q.Get("date_to")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportUsedParts(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "used-parts") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.UsedParts(r.Context(), optDate(q.Get("date_from")), optDate(q.Get("date_to")), optInt(q.Get("category_id")), optText(q.Get("brand")), optInt(q.Get("vehicle_id")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportVehicleMovement(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "vehicle-movement") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.VehicleMovement(r.Context(), optDate(q.Get("date_from")), optDate(q.Get("date_to")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportManagerSubordinates(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "manager-subordinates") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.ManagerSubordinates(r.Context(), optInt(q.Get("manager_id")))
	h.writeReport(w, rows, err)
}

func (h *Handler) reportRepairmanWorks(w http.ResponseWriter, r *http.Request) {
	if !h.requireReport(w, r, "repairman-works") {
		return
	}
	q := r.URL.Query()
	rows, err := h.reports.RepairmanWorks(r.Context(), optInt(q.Get("employee_id")), optDate(q.Get("date_from")), optDate(q.Get("date_to")), optInt(q.Get("vehicle_id")))
	h.writeReport(w, rows, err)
}

func (h *Handler) writeReport(w http.ResponseWriter, rows []map[string]any, err error) {
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func optText(v string) any {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return v
}

func optInt(v string) any {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil
	}
	return n
}

func optDate(v string) any {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil
	}
	return t
}
