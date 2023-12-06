package network

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(file)))

	if req.Method == http.MethodHead {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.WriteHeader(http.StatusOK)
	} else if req.Method == http.MethodGet {
		http.ServeFile(w, req, file)
	}
	return
}
