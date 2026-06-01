package constants

import "os"

const (
	Protocol = "http"
	Host     = "localhost"
	Port     = "7777" // should be the same as the environment variables of DOCKER_GIN_PORT and GIN_PORT
)

const (
	DevelopmentNamespace = "development"
	ProductionNamespace  = ""
	TestNamespace        = "test"
)

var (
	APIGroupBase       = getEnv("API_BASE_PATH", "api")
	DevelopmentBaseURL = APIGroupBase + "/" + DevelopmentNamespace + "/" + DevelopmentVersion
	ProductionBaseURL  = APIGroupBase + "/" + ProductionNamespace + "/" + ProductionVersion
	TestBaseURL        = APIGroupBase + "/" + TestNamespace + "/" + TestVersion
	CurrentBaseURL     = DevelopmentBaseURL
)

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var URLWhiteList = []string{
	"http",
	"https",
	"mailto",
	"tel",
	"ws",
}

var URLBlackList = []string{
	"javascript",
	"vbscript",
	"file",
	"data",
}
