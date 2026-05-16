package httpapi

import (
	"net/http"

	"autopark/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	crud    *application.CRUDService
	reports *application.ReportService
	auth    *application.AuthService
}

func NewRouter(crud *application.CRUDService, reports *application.ReportService, auth *application.AuthService) http.Handler {
	h := &Handler{crud: crud, reports: reports, auth: auth}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", h.login)
		r.Group(func(r chi.Router) {
			r.Use(h.authRequired)
			r.Get("/auth/me", h.me)
			r.Get("/meta/entities", h.entities)
			r.Get("/reports", h.reportList)
			r.Get("/reports/autopark", h.reportAutopark)
			r.Get("/reports/drivers", h.reportDrivers)
			r.Get("/reports/driver-distribution", h.reportDriverDistribution)
			r.Get("/reports/passenger-routes", h.reportPassengerRoutes)
			r.Get("/reports/mileage", h.reportMileage)
			r.Get("/reports/repairs-summary", h.reportRepairsSummary)
			r.Get("/reports/staff-hierarchy", h.reportStaffHierarchy)
			r.Get("/reports/garage", h.reportGarage)
			r.Get("/reports/vehicle-distribution", h.reportVehicleDistribution)
			r.Get("/reports/freight", h.reportFreight)
			r.Get("/reports/used-parts", h.reportUsedParts)
			r.Get("/reports/vehicle-movement", h.reportVehicleMovement)
			r.Get("/reports/manager-subordinates", h.reportManagerSubordinates)
			r.Get("/reports/repairman-works", h.reportRepairmanWorks)
			r.Route("/{entity}", func(r chi.Router) {
				r.Get("/", h.list)
				r.Post("/", h.create)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", h.get)
					r.Put("/", h.update)
					r.Patch("/", h.update)
					r.Delete("/", h.delete)
				})
			})
		})
	})

	return r
}

// Cross-Origin Resource Sharing.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
