package utils

import "time"

func InvokeFunctionWithInterval(duration time.Duration, functionToInvoke func()) {
	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		functionToInvoke()
	}
}
