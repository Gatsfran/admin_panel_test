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

func (r *Router) RegisterClientOrderRoutes() {
	r.router.HandleFunc("/api/v1/client_orders", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			var clientOrder entity.ClientOrder
			if err := json.NewDecoder(req.Body).Decode(&clientOrder); err != nil {
				log.Printf("Ошибка при декодировании тела запроса: %v", err)
				http.Error(w, "Неверный формат данных", http.StatusBadRequest)
				return
			}

			if err := clientOrder.Validate(); err != nil {
				log.Printf("Ошибка валидации заявки: %v", err)
				http.Error(w, "Невалидная заявка", http.StatusBadRequest)
			}

			if err := clientOrder.SetContactType(); err != nil {
				log.Printf("Ошибка при определении типа контакта: %v", err)
				http.Error(w, "Неверный формат контакта", http.StatusBadRequest)
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
	}).Methods("POST")

	protected := r.router.PathPrefix("/api/v1/admin").Subrouter()
	protected.Use(AuthMiddleware(r.cfg.JWTSecret))

	protected.HandleFunc("/client_orders", func(w http.ResponseWriter, req *http.Request) {
		clientOrders, err := r.db.ListClientOrder(req.Context())
		if err != nil {
			log.Printf("Ошибка при получении списка заказов: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(clientOrders); err != nil {
			log.Printf("Ошибка при кодировании списка заказов в JSON: %v", err)
			http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	protected.HandleFunc("/client_orders/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Printf("Неверный ID заказа: %v", err)
			http.Error(w, "Неверный ID заказа", http.StatusBadRequest)
			return
		}

		if err := r.db.DeleteClientOrder(req.Context(), id); err != nil {
			log.Printf("Ошибка при удалении заказа: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Заказ с ID %d удален", id)
	}).Methods("DELETE")
}
