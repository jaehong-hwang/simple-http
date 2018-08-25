package test

import shttp "github.com/jaehong-hwang/simple-http"

// GetEnv of Server
func GetEnv() shttp.ServerEnv {
	return shttp.ServerEnv{
		Port: 7857,
	}
}
