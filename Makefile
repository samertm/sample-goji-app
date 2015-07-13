.PHONY: serve watch-serve db-reset test docker-deps docker-build docker-run deploy-deps deploy docker check-to

serve:
	go install github.com/samertm/sample-gosu-app
	sample-gosu-app

watch-serve:
	$(shell while true; do $(MAKE) serve & PID=$$! ; echo $$PID ; inotifywait --exclude ".git" -r -e close_write . ; kill $$PID ; done)

db-reset:
	psql -h localhost -U sga -c "drop schema public cascade"
	psql -h localhost -U sga -c "create schema public"

test:
	go test $(ARGS) ./...

docker-deps:
	$(MAKE) -C postgres-docker docker-build
	$(MAKE) -C postgres-docker run-prod

docker-build:
	docker build -t sga .

docker-run:
	docker start sga-db # Did you run 'make docker-deps'?
	-docker top sga-app && docker rm -f sga-app
	docker run -d -p 8111:8000 --name sga-app --link sga-db:sga-db sga # Did you run 'make docker-build?'

docker: docker-build docker-run

# Must specify TO.
deploy-deps: check-to
	rsync -azP . $(TO):~/sample-gosu-app
	ssh $(TO) 'cd ~/sample-gosu-app && make docker-deps'

# Must specify TO.
deploy: check-to
	rsync -azP . $(TO):~/sample-gosu-app
	ssh $(TO) 'cd ~/sample-gosu-app && make docker'

check-to:
	ifndef TO
	    $(error TO is undefined)
	endif
