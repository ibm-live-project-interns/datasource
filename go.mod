module github.com/ibm-live-project-interns/datasource

go 1.23.0

require (
	github.com/ibm-live-project-interns/ingestor/shared v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.27.0 // indirect
	gorm.io/gorm v1.31.1 // indirect
)

replace github.com/ibm-live-project-interns/ingestor/shared => ../ingestor/shared
