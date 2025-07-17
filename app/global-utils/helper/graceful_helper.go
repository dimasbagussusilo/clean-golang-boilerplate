package helper

import "os"

var (
	ConsumerCount   = 0
	ExitAMQP        = make(chan os.Signal, 1)
	ExitKafka       = make(chan os.Signal, 1)
	ExitGrpcClient  = make(chan bool)
	ExitConsumer    = make(chan bool)
	ExitHTTP        = make(chan bool)
	ExitConcurrency = make(chan bool)
)
