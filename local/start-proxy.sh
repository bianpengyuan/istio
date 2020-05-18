killall -9 envoy | true

/home/bianpengyuan_google_com/workspace/go/src/istio.io/istio/out/linux_amd64/release/envoy \
-c /home/bianpengyuan_google_com/workspace/go/src/istio.io/istio/local/config.json \
--restart-epoch 0 --drain-time-s 45 --parent-shutdown-time-s 60 \
--max-obj-name-len 189 --local-address-ip-version v4 \
-l info --concurrency 2 \
&>/dev/null &