.PHONY: deploy docker

noopadapter: main.go
	go build -o noopbackend .
docker: noopbackend Dockerfile
	docker build . -t gcr.io/bpy-istio/noopadapter:head

deploy:
	kubectl -n istio-system apply -f noopadapter.yaml
