.PHONY: run-all run-api-gateway run-wiki-service run-entry-service run-comment-service run-version-service run-media-service clean

run-wiki-service:
	@echo "Running wiki-service..."
	cd wiki && ( \
		go build -o wiki-service . && \
		./wiki-service > wiki-service.log 2>&1 & \
		echo $$! > wiki-service.pid \
	)

run-entry-service:
	@echo "Running entry-service..."
	cd entry && ( \
		go build -o entry-service . && \
		./entry-service > entry-service.log 2>&1 & \
		echo $$! > entry-service.pid \
	)

run-comment-service:
	@echo "Running comment-service..."
	cd comment && ( \
		go build -o comment-service . && \
		./comment-service > comment-service.log 2>&1 & \
		echo $$! > comment-service.pid \
	)

run-version-service:
	@echo "Running version-service..."
	cd version && ( \
		go build -o version-service . && \
		./version-service > version-service.log 2>&1 & \
		echo $$! > version-service.pid \
	)

run-media-service:
	@echo "Running media-service..."
	cd media && ( \
		go build -o media-service . && \
		./media-service > media-service.log 2>&1 & \
		echo $$! > media-service.pid \
	)

run-api-gateway: 
	@echo "Running api-gateway..."
	cd gateway && ( \
		go build -o api-gateway . && \
		./api-gateway > api-gateway.log 2>&1 & \
		echo $$! > api-gateway.pid \
	)

run-auth-service: 
	@echo "Running auth-service..."
	cd auth && ( \
		go build -o auth-service . && \
		./auth-service > auth-service.log 2>&1 & \
		echo $$! > auth-service.pid \
	)

run-all: run-wiki-service run-entry-service run-comment-service run-version-service run-media-service run-api-gateway

clean:
	@echo "Stopping services..."
	-@cd wiki && kill `cat wiki-service.pid` && rm wiki-service.pid || true
	-@cd entry && kill `cat entry-service.pid` && rm entry-service.pid || true
	-@cd comment && kill `cat comment-service.pid` && rm comment-service.pid || true
	-@cd version && kill `cat version-service.pid` && rm version-service.pid || true
	-@cd media && kill `cat media-service.pid` && rm media-service.pid || true
	-@cd auth && kill `cat auth-service.pid` && rm auth-service.pid || true
	-@cd gateway && kill `cat api-gateway.pid` && rm api-gateway.pid || true
