module github.com/ibm-live-project-interns/datasource

go 1.23

require (
	github.com/ibm-live-project-interns/ingestor/shared v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
)

replace github.com/ibm-live-project-interns/ingestor/shared => ../ingestor/shared
