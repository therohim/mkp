# Makefile

# help to create new module seamlessly
# usage: make new-module name=moduleName
new-module:
	@mkdir -p src/$(name)
	@mkdir -p src/$(name)/controller
	@mkdir -p src/$(name)/entity
	@mkdir -p src/$(name)/enum
	@mkdir -p src/$(name)/model
	@mkdir -p src/$(name)/repository
	@mkdir -p src/$(name)/service
	@mkdir -p src/$(name)/validation
	@echo "Module $(name) created successfully."
	
server-start:
	nodemon --exec go run main.go --signal SIGTERM