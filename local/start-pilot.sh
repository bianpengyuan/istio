killall -9 pilot-discovery | true

/home/bianpengyuan_google_com/workspace/go/src/istio.io/istio/out/linux_amd64/pilot-discovery discovery \
--log_output_level=default:info \
--monitoringAddr=:15014 \
--keepaliveMaxServerConnectionAge=30m \
&>/dev/null &
