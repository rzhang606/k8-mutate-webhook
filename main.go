package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

const tlsConnection = true

// Webhook Server parameters
type whSvrParameters struct {
	port     int    // webhook server port
	certFile string // path to the x509 certificate for https
	keyFile  string // path to the x509 private key matching `CertFile`
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {

	// Ready body/request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
	}

	//mutate request
	mutated, err := mutate(body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
	}

	//Respond
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world\n")
}

func main() {
	var parameters whSvrParameters

	fmt.Print("Starting server ...")

	// get command line parameters
	flag.IntVar(&parameters.port, "port", 8443, "Webhook server port.")
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	var whsvr = &http.Server{
		Addr: fmt.Sprintf(":%v", parameters.port),
	}
	if tlsConnection == true {
		pair, err := tls.LoadX509KeyPair(parameters.certFile, parameters.keyFile)
		if err != nil {
			fmt.Printf("Failed to load key pair: %v", err)
		}

		whsvr = &http.Server{
			Addr:      fmt.Sprintf(":%v", parameters.port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		}
	}

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", mutateHandler)
	mux.HandleFunc("/", mainHandler)
	whsvr.Handler = mux

	// start webhook server in new rountine
	go func() {
		if tlsConnection == true {
			if err := whsvr.ListenAndServeTLS("", ""); err != nil {
				fmt.Printf("Failed to listen and serve webhook server: %v", err)
			}
		} else {
			if err := whsvr.ListenAndServe(); err != nil {
				fmt.Printf("Failed to serve webhook server %v", err)
			}
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
	whsvr.Shutdown(context.Background())
}
