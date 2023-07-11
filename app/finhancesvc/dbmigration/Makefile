create:
	migrate create -ext sql -dir sql $(name)

up:
	migrate -path sql -database ${FINHANCESVC_DBURI} up $(n)

down:
	migrate -path sql -database ${FINHANCESVC_DBURI} down $(n)

goto:
	migrate -path sql -database ${FINHANCESVC_DBURI} goto $(version)

force:
	migrate -path sql -database ${FINHANCESVC_DBURI} force $(version)
