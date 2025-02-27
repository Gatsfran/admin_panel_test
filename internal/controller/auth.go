package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateJWT(username string, jwtSecret string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string, jwtSecret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("недействительный токен")
}

func AuthMiddleware(jwtSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Запрос к защищённому маршруту: %s %s", r.Method, r.URL.Path)

			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				log.Printf("Токен отсутствует в запросе к %s", r.URL.Path)
				http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
				return
			}

			if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
				tokenString = tokenString[7:]
			}

			claims, err := ValidateJWT(tokenString, jwtSecret)
			if err != nil {
				log.Printf("Недействительный токен в запросе к %s: %v", r.URL.Path, err)
				http.Error(w, "Недействительный токен", http.StatusUnauthorized)
				return
			}

			log.Printf("Пользователь %s авторизован для запроса к %s", claims.Username, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
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
			log.Printf("Ошибка при получении хеша пароля для пользователя %s: %v", loginRequest.Username, err)
			http.Error(w, "Неверно указан логин или пароль", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginRequest.Password)); err != nil {
			http.Error(w, "Неверно указан логин или пароль", http.StatusUnauthorized)
			return
		}

		token, err := generateJWT(loginRequest.Username, r.cfg.JWTSecret)
		if err != nil {
			http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}).Methods("POST","OPTIONS")
}
