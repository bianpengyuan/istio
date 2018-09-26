.PHONY: deploy docker

noopadapter: main.go
	../../../../bin/gobuild.sh noopadapter istio.io/istio/pkg/version .
docker: noopadapter Dockerfile
	docker build . -t gcr.io/bpy-istio/noopadapter:head

deploy:
	kubectl -n istio-system apply -f noopadapter.yaml
