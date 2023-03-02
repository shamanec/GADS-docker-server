package android_server

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var streamImageChan = make(chan image.Image, 1)

var lastImageArray []byte
var dummyImage image.Image

func ConnectWS() {
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

		lastImageArray = message

		select {
		case streamImageChan <- dummyImage:
		default:
		}
	}
	c.Close()
}

type GadsStreamHandler struct {
	Next func() (image.Image, error)
}

func (h GadsStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	boundary := "\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"
	for {
		// get handler new image from imageChan
		_, err := h.Next()
		if err != nil {
			return
		}

		n, err := io.WriteString(w, boundary)
		if err != nil || n != len(boundary) {
			return
		}

		n, err = io.WriteString(w, string(lastImageArray))
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
		Next: func() (image.Image, error) {
			return <-streamImageChan, nil
		},
	}

	return &stream
}
