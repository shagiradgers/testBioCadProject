allBuild: buildApi buildApp

runApi:
	go run cmd/api/api.go
buildApi:
	go build cmd/api/api.go

runApp:
	go run cmd/app/app.go
buildApp:
	go build cmd/app/app.go