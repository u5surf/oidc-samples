COMMAND="./sample --client-id=${CLIENT_ID} --client-secret=${CLIENT_SECRET}"

if [ -n "$PROXY_HOST" ]; then COMMAND="${COMMAND} --proxy-host=${PROXY_HOST}"; fi
if [ -n "$PROXY_PORT" ]; then COMMAND="${COMMAND} --proxy-port=${PROXY_PORT}"; fi
if [ -n "$HOST" ]; then COMMAND="${COMMAND} --host=${HOST}"; fi
if [ -n "$PORT" ]; then COMMAND="${COMMAND} --port=${PORT}"; fi
if [ -n "$ISSUER" ]; then COMMAND="${COMMAND} --issuer=${ISSUER}"; fi
if [ -n "$REDIRECT_URL" ]; then COMMAND="${COMMAND} --redirect-url=${REDIRECT_URL}"; fi

eval $COMMAND
