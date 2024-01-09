bd:
	docker-compose -f docker-compose.dev.yaml build
bp:
	docker-compose build --no-cache
pup:
	docker-compose up -d
dup:
	docker-compose -f docker-compose.dev.yaml up -d
ld:
	docker-compose -f docker-compose.dev.yaml logs -f
lp:
	docker-compose logs -f
down-dev:
	docker-compose -f docker-compose.dev.yaml down
dd:
	docker-compose -f docker-compose.dev.yaml down
dv:
	docker-compose -f docker-compose.dev.yaml down -v --remove-orphans
dp:
	docker-compose down -v --remove-orphans
bup:
	docker-compose -f docker-compose.dev.yaml up --build -d --force-recreate --remove-orphans
exec-dev:
	docker exec -it go_bot sh

