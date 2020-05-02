package cameraDriver

import (
	"encoding/base64"
	"github.com/webcam"
)

var Camera = &Driver{}

type Driver struct {
	Camera 			*webcam.Webcam
	Status 			chan bool
	ReadFrame 		[]byte
	ReadImage 		string
}

func StartDevice(Port string) (*Driver, error) {

	camera, err := webcam.Open(Port)

	formatDesc := camera.GetSupportedFormats()
	var formats []webcam.PixelFormat
	for f := range formatDesc {
		formats = append(formats, f)
		camera.SetImageFormat(f, 1280, 720)
	}

	err = camera.StartStreaming()

	Camera.Camera = camera

	go func() {
		for {
			select {
			case <-Camera.Status:
				return
			default:
				frame, err := Camera.Camera.ReadFrame()
				if err != nil {
					continue
				}

				if len(frame) != 0 {
					Camera.ReadFrame = frame
					Camera.ReadImage = base64.StdEncoding.EncodeToString(frame)
				}
			}
		}
	}()

	return Camera, err
}