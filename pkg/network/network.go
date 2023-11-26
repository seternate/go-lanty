package network

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seternate/go-lanty/pkg/filesystem"
)

func GetOutboundIP() (ip net.IP, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	ip = conn.LocalAddr().(*net.UDPAddr).IP
	return
}

func ServeFileData(file string, w http.ResponseWriter, req *http.Request) (err error) {
	contentType, _ := filesystem.DetectContentTypeOfFile(file)

	w.Header().Set("Content-Type", contentType)
	builder := strings.Builder{}
	builder.WriteString("attachment; filename=")
	builder.WriteString(filepath.Base(file))
	w.Header().Set("Content-Disposition", builder.String())

	if req.Method == http.MethodHead {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		w.WriteHeader(http.StatusOK)
	} else if req.Method == http.MethodGet {
		http.ServeFile(w, req, file)
	}
	return
}
