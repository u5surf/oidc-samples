COMMAND="./integration-tests --client-id=${CLIENT_ID} --client-secret=${CLIENT_SECRET}" 

if [ -n "$SAMPLE_URL" ]; then COMMAND="${COMMAND} --sample-url=${SAMPLE_URL}"; fi
if [ -n "$API_URL" ]; then COMMAND="${COMMAND} --api-url=${API_URL}"; fi
if [ -n "$REDIRECT_URL" ]; then COMMAND="${COMMAND} --redirect-url=${REDIRECT_URL}"; fi

eval $COMMAND
