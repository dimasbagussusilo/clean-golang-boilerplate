#bin/bash
if [ "$HANDLER_TYPE" == "http" ]; then
    docker rm -f salesman-service 2> /dev/null || true
elif [ "$HANDLER_TYPE" == "consumer" ]; then
    docker rm -f consumer_$CONSUMER_NAME_ENV 2> /dev/null || true
fi