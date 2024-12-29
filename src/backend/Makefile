.PHONY: all build push $(SERVICES:%=build-%) run-all run-api-gateway run-wiki-service run-entry-service run-comment-service run-version-service run-media-service clean combine-swagger
# Define your services
SERVICES = auth comment entry gateway media translation wiki

# Default target when you run 'make' without arguments
all: build push

# Build all services
build: $(SERVICES:%=build-%)

# Push all services
push: $(SERVICES:%=push-%)

# Build each service
build-%:
	cd $* && sudo docker build -t klnstprx/$*-service .

# Push each service
push-%:
	sudo docker push klnstprx/$*-service

run-all: clean run-wiki-service run-entry-service run-comment-service run-version-service run-media-service run-api-gateway

clean:
	@echo "Stopping services..."
	-@cd wiki && test -e wiki-service.pid && kill `cat wiki-service.pid` 2>/dev/null && rm wiki-service.pid || true
	-@cd entry && test -e entry-service.pid && kill `cat entry-service.pid` 2>/dev/null && rm entry-service.pid || true
	-@cd comment && test -e comment-service.pid && kill `cat comment-service.pid` 2>/dev/null && rm comment-service.pid || true
	-@cd version && test -e version-service.pid && kill `cat version-service.pid` 2>/dev/null && rm version-service.pid || true
	-@cd media && test -e media-service.pid && kill `cat media-service.pid` 2>/dev/null && rm media-service.pid || true
	-@cd auth && test -e auth-service.pid && kill `cat auth-service.pid` 2>/dev/null && rm auth-service.pid || true
	-@cd gateway && test -e api-gateway.pid && kill `cat api-gateway.pid` 2>/dev/null && rm api-gateway.pid || true

combine-swagger:
	@echo "Generating Swagger documentation for all services..."
	cd wiki && swag init
	cd entry && swag init
	cd comment && swag init
	cd version && swag init
	cd media && swag init
	cd auth && swag init
	cd gateway && swag init
	@echo "Combining Swagger documentation..."
	npx swagger-combine swagger-config.json -o gateway/docs/swagger.json
	@echo "Swagger documentation combined successfully."&& sudo docker build -t klnstprx/$*-service .


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

run-translation-service: 
	@echo "Running translation-service..."
	cd translation && ( \
		go build -o translation-ser
