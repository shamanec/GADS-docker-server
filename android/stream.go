package android_server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var streamBytesChan = make(chan []byte, 1)

func ConnectGadsStreamWS() {
	u := url.URL{Scheme: "ws", Host: "localhost:1313", Path: ""}
	fmt.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		log.Fatal("dial:", err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		select {
		case streamBytesChan <- message:
		default:
		}
	}
	c.Close()
}

type GadsStreamHandler struct {
	Next func() ([]byte, error)
}

func (h GadsStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	boundary := "\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"
	for {
		// get handler new image from imageChan
		imageBytes, err := h.Next()
		if err != nil {
			return
		}

		n, err := io.WriteString(w, boundary)
		if err != nil || n != len(boundary) {
			return
		}

		n, err = io.WriteString(w, string(imageBytes))
		if err != nil {
			return
		}

		n, err = io.WriteString(w, "\r\n")
		if err != nil || n != 2 {
			return
		}
	}
}

func JpegStreamHandler() *GadsStreamHandler {
	// for each new image in imageChan update the handler
	stream := GadsStreamHandler{
		Next: func() ([]byte, error) {
			return <-streamBytesChan, nil
		},
	}

	return &stream
}
