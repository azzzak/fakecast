package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/azzzak/fakecast/api"
	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
)

var version string

func main() {
	var (
		host       string = ""
		root       string = "/fakecast"
		credential string = ""
		listenPort int    = 80
	)

	flag.StringVar(&host, "host", lookupEnvOrString("HOST", host), "host url")
	flag.StringVar(&root, "root", lookupEnvOrString("ROOT", root), "root of content directory")
	flag.StringVar(&credential, "credential", lookupEnvOrString("CREDENTIAL", credential), "access credential")
	flag.IntVar(&listenPort, "port", lookupEnvOrInt("PORT", listenPort), "port")

	flag.Parse()

	if host == "" {
		fmt.Println("You must set HOST env variable to proper work of app")
		os.Exit(1)
	}

	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}
	host = strings.TrimSuffix(host, "/")

	s, err := store.NewStore(root)
	if err != nil {
		fmt.Printf("Error while connecting to DB: %s\n", err)
		os.Exit(1)
	}
	defer s.Close()

	fs := fs.NewRoot(root)

	cfg := &api.Cfg{
		Store:      s,
		FS:         fs,
		Host:       host,
		Credential: credential,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", listenPort),
		Handler:      api.InitHandlers(cfg),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on port %d: %v\n", listenPort, err)
		}
	}()

	fmt.Printf("fakecast %s is running\n", version)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("fakecast is stopped")
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Env[%s]: %v", key, err)
			os.Exit(1)
		}
		return v
	}
	return defaultVal
}
