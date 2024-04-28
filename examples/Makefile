generate-schema:
	cd ../cmd &&  go run main.go > ../examples/ora-model.json

dbml:
	octopus-db-tools export dbml --input ora-model.json --output ora-model.dbml

gorm:
	octopus-db-tools generate gorm --input ora-model.json --output ora-model.go --package model	

plantuml:
	octopus-db-tools generate plantuml --input ora-model.json --output ora-model.plantuml 

.PHONY: generate-schema dbml gorm plantuml 