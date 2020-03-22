COMMAND="proxy"

if [ -n "$VERBOSE" ]; then COMMAND="${COMMAND} -v"; fi
if [ -n "$HOST" ]; then COMMAND="${COMMAND} --host=${HOST}"; fi
if [ -n "$PORT" ]; then COMMAND="${COMMAND} --port=${PORT}"; fi

eval $COMMAND
