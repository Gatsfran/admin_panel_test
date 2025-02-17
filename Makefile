# Старт контейнеров
run:
	docker-compose up -d --force-recreate
# Удаление контейнеров с принудительной их остановкой
down:
	docker-compose down
# Остановка контейнеров
stop:
	docker-compose stop
