package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
	"github.com/seternate/go-lanty/pkg/logging"
	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/router"
	"github.com/seternate/go-lanty/pkg/setting"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

func main() {
	signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignalCtx()
	errgrp, errCtx := errgroup.WithContext(signalCtx)

	settings, err := setting.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load settings")
	}
	logconfig := logging.Config{ConsoleLoggingEnabled: true}

	parseConfigFlags(&logconfig)
	parseServerFlags(&settings)
	printVersion := flag.Bool("version", false, "Prints the version information")
	flag.Parse()
	//Set default for settings.Host - bare LAN IP
	if len(settings.Host) == 0 {
		ip, err := network.GetOutboundIP()
		if err == nil {
			settings.Host = ip.String()
		}
	}

	log.Logger = logging.Configure(logconfig)

	log.Debug().Interface("logconfig", logconfig).Msg("logging configuration")
	log.Debug().Interface("settings", settings).Msg("settings")

	if *printVersion {
		fmt.Printf("%s - %s", setting.VERSION, runtime.Version())
		os.Exit(0)
	}

	err = updateClientServerURL(settings)
	if err != nil {
		log.Warn().Err(err).Msg("error during update of client \"settings.yaml\" serverurl")
	}

	handler := handler.NewHandler(&settings).
		WithGamehandler().
		WithUserhandler(errCtx, errgrp).
		WithDownloadHandler().
		WithChatHandler().
		WithFileHandler()
	log.Trace().Msg("handler created")
	gameRoutes := router.GameRoutes(handler)
	userRoutes := router.UserRoutes(handler)
	downloadRoutes := router.DownloadRoutes(handler)
	chatRoutes := router.ChatRoutes(handler)
	fileRoutes := router.FileRoutes(handler)
	log.Trace().Msg("routes created")
	router := router.NewRouter().
		WithRoutes(gameRoutes).
		WithRoutes(userRoutes).
		WithRoutes(downloadRoutes).
		WithRoutes(chatRoutes).
		WithRoutes(fileRoutes)
	log.Trace().Msg("router created")

	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(settings.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create listener")
	}
	httpServer := &http.Server{
		Handler: router,
		BaseContext: func(_ net.Listener) context.Context {
			return signalCtx
		},
	}

	var shutdownErr error
	errgrp.Go(func() error {
		log.Info().Msg("starting http server: " + listener.Addr().String())
		return httpServer.Serve(listener)
	})
	errgrp.Go(func() error {
		<-errCtx.Done()
		log.Info().Msgf("graceful shutdown of http server (%ds)", settings.GracefulShutdown)
		ctx, ctxCancel := context.WithTimeout(context.Background(), time.Duration(settings.GracefulShutdown)*time.Second)
		defer ctxCancel()
		shutdownErr = httpServer.Shutdown(ctx)
		return shutdownErr
	})

	err = errgrp.Wait()
	if err != nil || shutdownErr != nil {
		if err == http.ErrServerClosed && shutdownErr != nil {
			log.Fatal().Err(shutdownErr).Msg("server closed forcefully")
		} else if err != http.ErrServerClosed && err != context.Canceled && shutdownErr == nil {
			log.Fatal().Err(err).Msg("unexpected http server error")
		}
		log.Info().Msg("server closed")
	}
}

func parseConfigFlags(config *logging.Config) {
	flag.StringVar(&config.LogLevel, "loglevel", "info", "Sets the log level")
	flag.BoolVar(&config.FileLoggingEnabled, "logenablefile", false, "Enables logging to file")
	flag.StringVar(&config.Filename, "logfile", "lanty.log", "Sets the log filename")
	flag.StringVar(&config.Directory, "logdir", "log", "Sets the log directory")
	flag.IntVar(&config.MaxBackups, "logbackups", 0, "Sets the number of old logs to remain")
	flag.IntVar(&config.MaxSize, "logfilesize", 10, "Sets the size of the logs before rotating to new file")
	flag.IntVar(&config.MaxAge, "logage", 0, "Sets the maximum number of days to retain old logs")
}

func parseServerFlags(settings *setting.Settings) {
	flag.StringVar(&settings.Host, "host", settings.Host, "Hostname of the http server")
	flag.IntVar(&settings.Port, "port", settings.Port, "Port of the http server to listen on")
	flag.IntVar(&settings.GracefulShutdown, "graceful-shutdown", 10, "Timeout in seconds to wait for a graceful shutdown of the server")
	flag.StringVar(&settings.GameConfigDirectory, "game-config-dir", settings.GameConfigDirectory, "Directory of the game configuration files")
	flag.StringVar(&settings.GameFileDirectory, "game-file-dir", settings.GameFileDirectory, "Directory of the game files")
	flag.StringVar(&settings.GameIconDirectory, "game-icon-dir", settings.GameIconDirectory, "Directory of the game icons")
	flag.StringVar(&settings.FileUploadDirectory, "file-upload-dir", settings.FileUploadDirectory, "Directory of files uploaded by clients")
}

func updateClientServerURL(settings setting.Settings) (err error) {
	path := filepath.Join(setting.CLIENT_DOWNLOAD_DIRECTORY, setting.CLIENT_DOWNLOAD_FILE)
	readBuffer, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("error reading \"%s\": %w", path, err)
		return
	}

	file, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		err = fmt.Errorf("error opening \"%s\": %w", path, err)
		return
	}
	defer file.Close()

	reader := bytes.NewReader(readBuffer)
	zipReader, err := zip.NewReader(reader, int64(len(readBuffer)))
	if err != nil {
		err = fmt.Errorf("error creating zip reader \"%s\": %w", path, err)
		return
	}

	writeBuffer := bytes.Buffer{}
	zipWriter := zip.NewWriter(&writeBuffer)

	for _, zipFile := range zipReader.File {
		zipFileReader, err := zipFile.Open()
		if err != nil {
			err = fmt.Errorf("error opening zip file \"%s\": %w", zipFile.Name, err)
			return err
		}

		header, err := zip.FileInfoHeader(zipFile.FileInfo())
		if err != nil {
			err = fmt.Errorf("error creating zip file header \"%s\": %w", zipFile.Name, err)
			return err
		}

		header.Name = zipFile.Name
		header.Method = zip.Deflate
		zipFileWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			err = fmt.Errorf("error creating zip file writer \"%s\": %w", zipFile.Name, err)
			return err
		}

		if zipFile.Name == "settings.yaml" {
			zipFileBuffer, err := io.ReadAll(zipFileReader)
			if err != nil {
				err = fmt.Errorf("error reading settings file \"%s\": %w", zipFile.Name, err)
				return err
			}
			tmpSettingsFile := make(map[string]interface{})
			err = yaml.Unmarshal(zipFileBuffer, tmpSettingsFile)
			if err != nil {
				err = fmt.Errorf("error unmarshal settings file \"%s\": %w", zipFile.Name, err)
				return err
			}
			if _, found := tmpSettingsFile["serverurl"]; found {
				tmpSettingsFile["serverurl"] = fmt.Sprintf("http://%s:%d", settings.Host, settings.Port)
			} else {
				err = errors.New("settings file missing \"serverurl\" key")
				return err
			}
			tmpSettingsFileBuffer, err := yaml.Marshal(tmpSettingsFile)
			if err != nil {
				err = fmt.Errorf("error marshal settings file \"%s\": %w", zipFile.Name, err)
				return err
			}
			_, err = io.Copy(zipFileWriter, bytes.NewReader(tmpSettingsFileBuffer))
			if err != nil {
				err = fmt.Errorf("error writting to settings file buffer \"%s\": %w", zipFile.Name, err)
				return err
			}
		} else {
			_, err = io.Copy(zipFileWriter, zipFileReader)
			if err != nil {
				err = fmt.Errorf("error writting to zip file buffer \"%s\": %w", zipFile.Name, err)
				return err
			}
		}
		zipFileReader.Close()
	}
	err = zipWriter.Close()
	if err != nil {
		err = fmt.Errorf("error writting central directory to buffer: %w", err)
		return
	}

	_, err = io.Copy(file, &writeBuffer)

	return
}
