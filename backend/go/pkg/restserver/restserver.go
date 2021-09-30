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
	mm := make(map[string]string)
	mm["message"] = message
	out, err := json.Marshal(mm)
	if err != nil {
		return err
	}
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

	// enable cors
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	ch.logger.Printf("body is: '%s'", string(data))

	var clipperRequest clips.ClipsRequest
	err = json.Unmarshal(data, &clipperRequest)
	if err != nil {
		ch.logger.Printf("could not unmarshal request: %s", err)
	}

	err = ch.cs.CreateClips(clips.ToClip(&clipperRequest))

	if err != nil {
		ch.logger.Printf("could not create clip  %s", err)
	}

	writeJsonMessage(rw, "OK", 200)
}

func addHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//anyone can make a CORS request (not recommended in production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//only allow GET, POST, and OPTIONS
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		//Since I was building a REST API that returned JSON, I set the content type to JSON here.
		w.Header().Set("Content-Type", "application/json")
		//Allow requests to have the following headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")
		//if it's just an OPTIONS request, nothing other than the headers in the response is needed.
		//This is essential because you don't need to handle the OPTIONS requests in your handlers now
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

type ClipperRESTServer struct {
	logger *log.Logger
	server *http.Server
	sm     *http.ServeMux
}

func (crs ClipperRESTServer) Addr() string {
	return crs.server.Addr
}

func NewClipperRESTServer(logger *log.Logger, cs *clippersrvc.ClipperService, host, port string) *ClipperRESTServer {
	ch := ClipperHandler{logger, cs}
	sm := &http.ServeMux{}
	sm.Handle("/clip", addHeaders(&ch))
	server := &http.Server{Handler: sm, Addr: fmt.Sprintf("%s:%s", host, port)}
	return &ClipperRESTServer{logger: logger, server: server, sm: sm}
}

func (crs *ClipperRESTServer) Run() error {
	listener, err := net.Listen("tcp", crs.Addr())
	if err != nil {
		crs.logger.Fatalf("could no listen on %s: %s", crs.Addr(), err)
	}

	return crs.server.Serve(listener)
}

func (crs *ClipperRESTServer) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return crs.server.Shutdown(ctx)
}
