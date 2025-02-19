package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Gatsfran/admin_panel_test/internal/entity"
	"github.com/Gatsfran/admin_panel_test/internal/repo"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
	db     *repo.DB
}

func New(db *repo.DB) *Router {
	r := mux.NewRouter()
	router := &Router{
		router: r,
		db:     db,
	}
	router.registerClientOrderRoutes()
	return router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

var (
	jwtSecret = []byte("тут должен быть секретный ключ")
)
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), 
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("недействительный токен")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		log.Printf("Пользователь %s авторизован", claims.Username)
		next.ServeHTTP(w, r)
	})
}

func (r *Router) RegisterAuthRoutes() {
	r.router.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(req.Body).Decode(&loginRequest); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		passwordHash, err := r.db.GetPasswordHash(req.Context(), loginRequest.Username)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}

		if passwordHash != loginRequest.Password {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

		token, err := GenerateJWT(loginRequest.Username)
		if err != nil {
			http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}).Methods("POST")
}

func (r *Router) registerClientOrderRoutes() {
	r.router.HandleFunc("/client_orders", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			clientOrders, err := r.db.ListClientOrder(req.Context())
			if err != nil {
				log.Printf("Ошибка при получении списка заявок: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := r.db.CreateClientOrder(req.Context(), &clientOrder); err != nil {
				log.Printf("Ошибка при создании заявки: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(clientOrder); err != nil {
				log.Printf("Ошибка при кодировании заявки в JSON: %v", err)
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	protected := r.router.PathPrefix("/admin").Subrouter()
	protected.Use(AuthMiddleware)

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
				log.Printf("Ошибка при удалении заявок: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Заявка с ID %d удалена", id)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("DELETE")
}
