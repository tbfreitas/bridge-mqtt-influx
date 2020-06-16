package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/influxdata/influxdb/client/v2"
)

type Metric struct {
	Fieldname       string
	Fieldvalue      string
	Measurementname string
}

const (
	INFLUXDB_NAME = "teste_db"
	USERNAME      = "admin"
	PASSWORD      = "admin"
	HOST          = "http://influxdb:8086"
)

var hc = createClient()

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
	var metric Metric
	json.Unmarshal([]byte(msg.Payload()), &metric)
	fmt.Printf("Nome da metrica: %s, Valor do campo: %s, nome da métrica: %s", metric.Fieldname, metric.Fieldvalue, metric.Measurementname)
	bp := createMetrics(metric.Fieldname, metric.Fieldvalue, metric.Measurementname)
	hc.Write(bp)
}

func createMetrics(fn string, fv string, mn string) client.BatchPoints {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  INFLUXDB_NAME,
		Precision: "s",
	})

	tags := map[string]string{"cpu": "cpu-total"}

	fields := map[string]interface{}{
		fn: fv,
	}

	pt, _ := client.NewPoint(mn, tags, fields, time.Now())
	bp.AddPoint(pt)
	return bp
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := MQTT.NewClientOptions().AddBroker("mqtt://localhost:1883")
	opts.SetClientID("mac-tarcisio")
	opts.SetDefaultPublishHandler(f)
	topic := "/teste/memory"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	MQTTclient := MQTT.NewClient(opts)

	if token := MQTTclient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Conectado ao broker tcp://localhost:1883\n")
	}

	<-c
}

func createClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "admin",
		Password: "admin",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return c
}
