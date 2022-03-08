package environment

import (
	"os"
)

var envName = os.Getenv("ENV_NAME")

func decodeLocation() string {
	switch os.Getenv("LOCATION") {
	case "hk":
		return "ap-east-1"
	case "us":
		return "us-east-1"
	case "ie":
		return "eu-west-1"
	case "ap":
		return "ap-southeast-1"
	}
	return ""
}

var region = decodeLocation()
var cluster = os.Getenv("CLUSTER")

func IsProd() bool {
	return envName == "prod" || envName == "prd"
}

func EnvName() string {
	return envName
}

func Region() string {
	return region
}

func Cluster() string {
	return cluster
}

