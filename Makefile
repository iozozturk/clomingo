build:
	go build -o clomingo ./cmd/main.go
run:
	go run ./cmd
deploy:
	gcloud app deploy --project=clomingo-dev
deploy-no-promote:
	gcloud app deploy --project=clomingo-dev --no-promote
deploy-prod:
	gcloud app deploy prod-app.yaml --project=clomingo-prod
deploy-prod-no-promote:
	gcloud app deploy prod-app.yaml --project=clomingo-prod --no-promote