.PHONY: dc-up-8.0
dc-up-8.0:
	VERSION="8.0" docker compose up -d

.PHONY: dc-up-5.7
dc-up-5.7:
	VERSION="5.7" docker compose up -d

.PHONY: dc-down
dc-down:
	docker compose down

.PHONY: exec
exec:dc-up-8.0
	docker exec -it rdb /bin/bash

.PHONY: exec-mysql
exec-mysql:dc-up-8.0
	docker exec -it rdb /bin/bash -c "mysql -u root -ppassword rdb"

.PHONY: which
which:dc-up-8.0
	docker exec -i rdb sh /opt/mysql/which.sh

.PHONEY: init-db
init-db:dc-up-8.0
	docker exec -i rdb sh /opt/mysql/setup.sh
