1. make build & make docker.proxyv2
2. Dump the desired proxy config `kc exec $POD -n $POD_NAMESPACE curl localhost:15000/config_dump > ./local/config.yaml`
   Then edit several things: 
      1) Edit xds-grpc discovery address to localhost. (s/"address": "istiod.istio-system.svc"/"address": "localhost")
      2) Remove the following block:
         ```
          ,
          "tls_certificate_sds_secret_configs": [
           {
            "name": "default",
            "sds_config": {
             "api_config_source": {
              "api_type": "GRPC",
              "grpc_services": [
               {
                "envoy_grpc": {
                 "cluster_name": "sds-grpc"
                }
               }
              ]
             }
            }
           }
          ]
         ```
      3) Replace root cert path. (s/.\/var\/run\/secrets\/istio\/root-cert.pem//home/bianpengyuan_google_com/workspace/go/src/istio.io/istio/local/var/run/secrets/istio-dns/key.pem)
      4) Replace `stats matcher` to
      ```
      "stats_matcher": {
       "inclusion_list": {
        "patterns": [
         {
          "regex": ".*"
         }
        ]
       }
      }
      ```
3. run ./local/run.sh

localhost:15000 could have dump.

Also steps to run proxy locally with remote istiod from https://github.com/istio/istio/issues/23858:

port forward istiod to 15012

mkdir -p out/certs/vm
generate_cert --mode citadel \
    -out-cert out/certs/vm/cert-chain.pem -out-priv out/certs/vm/key.pem \
    -host spiffe://cluster.local/ns/vmtest/sa/default
kubectl get cm -n istio-system istio-ca-root-cert -ojsonpath='{.data.root-cert\.pem}' > out/certs/vm/root-cert.pem
OUTPUT_CERTS=./out/certs/vm PROV_CERT=./out/certs/vm CA_ADDR=localhost:15012 TERMINATION_DRAIN_DURATION_SECONDS=0 ISTIO_META_ISTIO_VERSION=1.9.0 PROXY_CONFIG="$(cat ~/kube/local/proxyconfig.yaml | envsubst)" go run ./pilot/cmd/pilot-agent proxy sidecar --templateFile ./tools/packaging/common/envoy_bootstrap_v2.json --proxyLogLevel=debug |& h info debug warn error critical
where proxyconfig.yaml is

binaryPath: $GOPATH/src/istio.io/istio/out/linux_amd64/debug/envoy
configPath: $HOME/kube/local/proxy
discoveryAddress: localhost:15012
statusPort: 15020
Switching debug image for release image and things work
