package main

import (
	"fmt"
	 MQTT "github.com/eclipse/paho.mqtt.golang"
	 "github.com/influxdata/influxdb/client/v2"
	 "os"
	 "log"
	 "math/rand"
	 "os/signal"
	 "syscall"
	 "time"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message){ 
   fmt.Printf("MSG: %s\n", msg.Payload())
	 httpClient := createClient()
	 createMetrics(httpClient)
}

func createMetrics(c client.Client) {

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "teste_db",
		Precision: "us",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	tags := map[string]string{
		"cluster": "host1",
		"host":    fmt.Sprintf("192.168.%d.%d", 1, rand.Intn(100)),
	}

	fields := map[string]interface{}{
		"cpu_usage":  rand.Float64() * 100.0,
		"disk_usage": rand.Float64() * 100.0,	
	}

	eventTime := time.Now().Add(time.Second * -20)

	point, err := client.NewPoint(
		"node_status",
		tags,
		fields,
		eventTime.Add(time.Second*10),
	)

	bp.AddPoint(point)
}
	
func createClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     "http://localhost:8086",
			Username: "admiaaan",
			Password: "admin",
	})

	if err != nil {
			log.Fatalln("Error: ", err)
	}
	return c
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

   opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
   opts.SetClientID("mac-tarcisio")
   opts.SetDefaultPublishHandler(f)
	 topic := "/teste/memory"

	 opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
						panic(token.Error())
			}
		}
		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
				panic(token.Error())
		} else {
				fmt.Printf("Connected to server\n")
		}

	<-c
} //en