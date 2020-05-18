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