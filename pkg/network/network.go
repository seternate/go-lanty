package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/seternate/go-lanty-server/pkg/filesystem"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func ServeFileData(f string, w http.ResponseWriter, req *http.Request) (err error) {
	contentType, err := filesystem.DetectContentTypeOfFile(f)
	if err != nil {
		err = errors.New(fmt.Sprintf("Can not detect content type of file '%s'", f))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(f)))
	if req.Method == http.MethodHead {
		fi, err := os.Stat(f)
		if err != nil {
			err = errors.New(fmt.Sprintf("Can not determine filesize of '%s'", f))
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))
		w.WriteHeader(200)
	} else if req.Method == http.MethodGet {
		http.ServeFile(w, req, f)
	}

	return nil
}
