package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/microservices/types"
)

const wsEndpoint = "ws://localhost:30000/ws"

var sendInterval = time.Second

func main() {

	obuids := generateOBUIDs(20)

	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)

	if err != nil {
		log.Fatal(err)
	}

	for {

		for i := 0; i < len(obuids); i++ {

			lat, long := genLocation()

			data := types.OBUData{
				OBUID: obuids[i],
				Lat:   lat,
				Long:  long,
			}

			// fmt.Printf("%+v\n", data)

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(sendInterval)
	}
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}
func genLocation() (float64, float64) {
	return genCoord(), genCoord()
}
func generateOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}
