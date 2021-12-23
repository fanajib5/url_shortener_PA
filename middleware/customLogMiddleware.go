package middleware

import (
	"fmt"
	"time"
)

func CustomLogMiddleware(method string, uri string, status int) {
	timestamp := time.Now().Format(time.RFC3339Nano)
	fmt.Printf("%s method=%s, uri=%s, status=%d __by CustomLogMiddleware__\n", timestamp, method, uri, status)
}
