package controllers

import (
	"../repository"
	"../validation"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strconv"
)

//When /hit is called
func HandleHit(w http.ResponseWriter, r *http.Request) {

	//parsing url params
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		showPixel(w)
		return
	}

	//Validation of the client id and visitor id GET params
	if !validator.ClientId(r.FormValue("cId")) || !validator.VisitorId(r.FormValue("vId")) {
		fmt.Println("Invalid client or visitor id " + r.FormValue("cId") + " " + r.FormValue("vId"))
		showPixel(w)
		return
	}

	//getting the client id in the correct type
	clientId64, err := strconv.ParseUint(r.FormValue("cId"), 10, 32)
	//ParseUint return only uint64, I need to convert it to uint32
	clientId32 := uint32(clientId64)

	//getting the VistorId id in the correct type
	visitorId64, err := strconv.ParseUint(r.FormValue("vId"), 10, 64)
	fmt.Println(clientId32, visitorId64)

	//checking the client status (existance) to make sure he's autorized
	//if trackingStatus and predictionStatus if false the client does not exit or at least no longer
	trackingStatus, predictionStatus := repository.GetClientsStatus(clientId32)
	if !trackingStatus && !predictionStatus {
		fmt.Println("Client id is not authorized " + r.FormValue("cId"))
		showPixel(w)
		return
	}

	//validation of the rest of the params and binding to the right structure
	//this function return a json string, a client id and a visitor id
	hitInJson, err := validator.GetHitType(r.FormValue("type"), r.Form)
	if err != nil {
		fmt.Println(err)
		showPixel(w)
		return
	}

	//fmt.Printf("%d -> %d -> %s\n", cId, vId, r.FormValue("type"))

	//if the javascript framework ask for prediction result
	if _, ok := r.Form["callback"]; ok {

		//getting the prediction status in redis
		//visitorBuckets := repository.GetVisitor(clientId32, visitorId64)

		visitorBuckets := "toto"

		//retruning it with javascript function
		showBuckets(w, visitorBuckets)
	} else {
		showPixel(w)
	}

	//sending the hit to real time prediction stream
	if predictionStatus {
		repository.SendToStream(visitorId64, hitInJson)
	}

	//inserting the hit in the long terme storage stream (firehose)
	if trackingStatus {
		//SendToFirehose create a buffer of hits and send them to the stream only when the buffer is full
		//repository.SendToFirehose(hitInJson)
	}

}

//to return an empty prediction
func showBlancResult(w http.ResponseWriter) {
	fmt.Fprintf(w, "StormizeHandleBucket({});")
}


//to return all predictions found for the visitor
func showBuckets(w http.ResponseWriter, buckets string) {
	fmt.Fprintf(w, "StormizeHandleBucket("+buckets+");")
}


//if the prediction is disabled and we only want to track visitor behavior we return a 1x1 pixel
func showPixel(w http.ResponseWriter) {

	//creating a transparent 1x1 pixel
	background := image.NewRGBA(image.Rect(0, 0, 1, 1))
	draw.Draw(background, background.Bounds(), image.Transparent, image.ZP, draw.Src)
	var img image.Image = background

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	//sending header for png image
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	//output
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
