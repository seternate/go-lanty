package swagger

import (
	"fmt"

	"github.com/seternate/go-lanty/docs"
)

func InitInfo(version string, host string, port int, basePath string, schemes []string) {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", host, port)
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = schemes
}
