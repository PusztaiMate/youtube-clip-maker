package restserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/PusztaiMate/clipper-go-backend/clips"
	"github.com/PusztaiMate/clipper-go-backend/pkg/clippersrvc"
)

type HTTPError struct {
	message   string
	errorCode int
}

func NewHTTPError(message string, errorCode int) *HTTPError {
	return &HTTPError{message: message, errorCode: errorCode}
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error - code: %d, msg: %s", h.errorCode, h.message)
}

func (h *HTTPError) writeAsError(rw http.ResponseWriter) error {
	err := writeJsonMessage(rw, h.message, h.errorCode)

	return err
}

func writeJsonMessage(rw http.ResponseWriter, message string, code int) error {
	out, err := json.Marshal(message)
	if err != nil {
		return err
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	_, err = rw.Write(out)
	if err != nil {
		return err
	}

	return nil
}

type ClipperHandler struct {
	logger *log.Logger
	cs     *clippersrvc.ClipperService
}

func (ch *ClipperHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var data []byte
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewHTTPError("could not open body", 400).writeAsError(rw)
	}

	ch.logger.Printf("body is: '%s'", string(data))

	var clipperRequest clips.ClipsRequest
	err = json.Unmarshal(data, &clipperRequest)
	if err != nil {
		ch.logger.Printf("could not unmarshal request: %s", err)
	}

	// err = ch.cs.CreateClips(&clipperRequest)
	ch.logger.Printf("creating clips based on: %#v", &clipperRequest)
	ch.logger.Println("clips present in request:")
	for i, c := range clipperRequest.GetClips() {
		ch.logger.Printf("%d: %#v", i, c)
	}
	if err != nil {
		ch.logger.Printf("could not create clip  %s", err)
	}
}

type ClipperRESTServer struct {
	logger     *log.Logger
	server     *http.Server
	sm         *http.ServeMux
	port, host string
	shutdown   chan struct{}
}

func (crs ClipperRESTServer) Addr() string {
	return fmt.Sprintf("%s:%s", crs.host, crs.port)
}

func NewClipperRESTServer(logger *log.Logger, cs *clippersrvc.ClipperService, host, port string) *ClipperRESTServer {
	shutdownCh := make(chan struct{})
	ch := ClipperHandler{logger, cs}
	sm := &http.ServeMux{}
	sm.Handle("/clip", &ch)
	server := &http.Server{Handler: sm}
	return &ClipperRESTServer{logger: logger, server: server, sm: sm, shutdown: shutdownCh, host: host, port: port}
}

func (crs *ClipperRESTServer) Run() context.CancelFunc {
	listener, err := net.Listen("tcp", crs.Addr())
	if err != nil {
		crs.logger.Fatalf("could no listen on %s: %s", crs.Addr(), err)
	}

	errChan := make(chan error)

	go func() {
		err = crs.server.Serve(listener)
		errChan <- err
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		select {
		case <-crs.shutdown:
			crs.server.Shutdown(ctx)
		case <-errChan:
			crs.logger.Fatalf("server failed with '%s'", err)
		}
	}()

	return cancel
}

func (crs *ClipperRESTServer) Shutdown() {
	crs.shutdown <- struct{}{}
}
