package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Gatsfran/admin_panel_test/internal/entity"
	"github.com/gorilla/mux"
)

func (r *Router) registerClientOrderRoutes() {
	r.router.HandleFunc("/admin/client_orders", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			clientOrders, err := r.db.ListClientOrder(req.Context())
			if err != nil {
				log.Printf("Ошибка при получении списка заявок: %v", err)
				http.Error(w, "Не удалось получить список заявок", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(clientOrders); err != nil {
				log.Printf("Ошибка при кодировании списка заявок в JSON: %v", err)
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}

		case http.MethodPost:
			var clientOrder entity.ClientOrder
			if err := json.NewDecoder(req.Body).Decode(&clientOrder); err != nil {
				log.Printf("Ошибка при декодировании тела запроса: %v", err)
				http.Error(w, "Неверный формат данных", http.StatusBadRequest)
				return
			}
			if err := r.db.CreateClientOrder(req.Context(), &clientOrder); err != nil {
				log.Printf("Ошибка при создании заявок: %v", err)
				http.Error(w, "Не удалось создать заявку", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(clientOrder); err != nil {
				log.Printf("Ошибка при кодировании заявок в JSON: %v", err)
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	protected := r.router.PathPrefix("/admin").Subrouter()
	protected.Use(authMiddleware(r.cfg.JWTSecret)) // Используем JWTSecret из конфигурации

	protected.HandleFunc("/client_orders/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Printf("Неверный ID заявки: %v", err)
			http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case http.MethodDelete:
			if err := r.db.DeleteClientOrder(req.Context(), id); err != nil {
				log.Printf("Ошибка при удалении заявки: %v", err)
				http.Error(w, "Не удалось удалить заявку", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Заявка с ID %d удален", id)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("DELETE")
}
