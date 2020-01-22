COMMAND="./sample --client-id=${CLIENT_ID} --client-secret=${CLIENT_SECRET}" 

[[ -n "$HOST" ]] && COMMAND="${COMMAND} --host=${HOST}"
[[ -n "$PORT" ]] && COMMAND="${COMMAND} --port=${PORT}"
[[ -n "$ISSUER" ]] && COMMAND="${COMMAND} --issuer=${ISSUER}"
[[ -n "$REDIRECT_URL" ]] && COMMAND="${COMMAND} --redirect-url=${REDIRECT_URL}"

eval $COMMAND
