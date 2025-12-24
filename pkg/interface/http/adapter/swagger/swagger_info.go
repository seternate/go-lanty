package swagger

import (
	"fmt"

	"github.com/seternate/go-lanty/docs"
)

func InitInfo(version string, basePath string, port int) {
	docs.SwaggerInfo.Title = "Lanty"
	docs.SwaggerInfo.Description = "Dummy description" //TODO: Replace with actual description
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.InfoInstanceName = ""
	docs.SwaggerInfo.Schemes = []string{"http"} //TODO: Replace with actual schemes
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port) //TODO: Replace with actual host
}
