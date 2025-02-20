package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func validateJWT(tokenString string, jwtSecret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("токен истек: %v", claims.ExpiresAt)
	}

	return claims, nil
}

func authMiddleware(jwtSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
				return
			}

			claims, err := validateJWT(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Недействительный токен", http.StatusUnauthorized)
				return
			}

			log.Printf("Пользователь %s авторизован", claims.Username)
			next.ServeHTTP(w, r)
		})
	}
}

func (r *Router) registerAuthRoutes() {
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
	}).Methods("POST")
}
