GPJC_RELEASE_URL := https://github.com/Ethernal-Tech/private-join-and-compute/releases/download/v0.1.0/bazel-bin.tar.gz
GPJC_MULTIPLE_RELEASE_URL := https://github.com/Ethernal-Tech/private-join-and-compute/releases/download/v0.1.0/bazel-bin-multiple.tar.gz
GPJC_RELEASE_FOLDER := private-join-and-compute

GPJC_API_RELEASE_ULR := https://github.com/Ethernal-Tech/gpjc-api/releases/download/v0.1.1/gpjc-api
GPJC_API_MULTIPLE_RELEASE_ULR := https://github.com/Ethernal-Tech/gpjc-api/releases/download/v0.1.3/gpjc-api-multiple

fetch-releases:
	@echo "Checking if $(GPJC_RELEASE_FOLDER) exists..."
	if [ ! -d "$(GPJC_RELEASE_FOLDER)" ]; then \
		echo "$(GPJC_RELEASE_FOLDER) does not exist. Creating..."; \
		mkdir -p "$(GPJC_RELEASE_FOLDER)"; \
	fi
	@echo "Fetching release from $(GPJC_RELEASE_URL)"
	curl -L $(GPJC_RELEASE_URL) -o /tmp/release.tar.gz
	tar -xzvf /tmp/release.tar.gz -C "$(GPJC_RELEASE_FOLDER)"
	rm /tmp/release.tar.gz
	@echo "Release $(GPJC_RELEASE_FOLDER) fetched successfully"
	@echo "Fetch API"
	curl -L $(GPJC_API_RELEASE_ULR) -o gpjc-api
	@echo "Give permissions to API exe"
	chmod +x gpjc-api

fetch-releases-multiple-machines:
	@echo "Checking if $(GPJC_RELEASE_FOLDER) exists..."
	if [ ! -d "$(GPJC_RELEASE_FOLDER)" ]; then \
		echo "$(GPJC_RELEASE_FOLDER) does not exist. Creating..."; \
		mkdir -p "$(GPJC_RELEASE_FOLDER)"; \
	fi
	@echo "Fetching release from $(GPJC_MULTIPLE_RELEASE_URL)"
	curl -L $(GPJC_MULTIPLE_RELEASE_URL) -o /tmp/release.tar.gz
	tar -xzvf /tmp/release.tar.gz -C "$(GPJC_RELEASE_FOLDER)"
	rm /tmp/release.tar.gz
	@echo "Release $(GPJC_RELEASE_FOLDER) fetched successfully"
	@echo "Fetch API"
	curl -L $(GPJC_API_MULTIPLE_RELEASE_ULR) -o gpjc-api
	@echo "Give permissions to API exe"
	chmod +x gpjc-api

create-certs: 
	chmod +x image/gpjc_scripts/ca_script.sh
	./image/gpjc_scripts/ca_script.sh

run-docker: create-certs
	docker-compose up --build -d

stop-docker: 
	docker-compose down --rmi local -v

test: run-docker
	sleep 90
	$(MAKE) -C playwright-tests test