package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/d2r2/go-dht"
	"github.com/d2r2/go-shell"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"

	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)

func main() {
	defer logger.FinalizeLogger()

	lg.Notify("***************************************************************************************************")
	lg.Notify("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	lg.Notify("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	lg.Notify("***************************************************************************************************")
	lg.Notify("*** Massive stress test of sensor reading, printing in the end summary statistical results")
	lg.Notify("***************************************************************************************************")
	// Uncomment/comment next line to suppress/increase verbosity of output
	logger.ChangePackageLogLevel("dht", logger.InfoLevel)

	// create context with cancellation possibility
	ctx, cancel := context.WithCancel(context.Background())
	// use done channel as a trigger to exit from signal waiting goroutine
	done := make(chan struct{})
	defer close(done)
	// build actual signal list to control
	signals := []os.Signal{os.Kill, os.Interrupt}
	if shell.IsLinuxMacOSFreeBSD() {
		signals = append(signals, syscall.SIGTERM)
	}
	// run goroutine waiting for OS termination events, including keyboard Ctrl+C
	shell.CloseContextOnSignals(cancel, done, signals...)

	sensorType := dht.DHT11
	// sensorType := dht.AM2302
	//sensorType := dht.DHT12
	pin, err := strconv.Atoi(os.Getenv("DHT11PIN"))
	if err != nil {
		lg.Error(err)
		lg.Info("exited")
		return
	}
	totalRetried := 0
	totalMeasured := 0
	totalFailed := 0
	term := false

	// connect to Mqtt broker
	cli := connectToMqtt()

	for {
		// Read DHT11 sensor data from specific pin, retrying 10 times in case of failure.
		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithContextAndRetry(ctx, sensorType, pin, false, 10)
		totalMeasured++
		totalRetried += retried
		if err != nil && ctx.Err() == nil {
			totalFailed++
			lg.Error(err)
			continue
		}
		// print temperature and humidity
		if ctx.Err() == nil {
			lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
				sensorType, temperature, humidity, retried)
		}

		// publish temperature status to mqtt broker
		publishToMqtt(cli, temperature, humidity)

		select {
		// Check for termination request.
		case <-ctx.Done():
			lg.Errorf("Termination pending: %s", ctx.Err())
			term = true
			// sleep 1.5-2 sec before next round
			// (recommended by specification as "collecting period")
		case <-time.After(2000 * time.Millisecond):
		}
		if term {
			break
		}
	}
	lg.Info("exited")
}

func connectToMqtt() *client.Client {
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("receive-client"),
	})
	if err != nil {
		panic(err)
	}
	return cli
}

func publishToMqtt(cli *client.Client, temperature float32, humidity float32) {
	deviceTwinUpdate := "$hw/events/device/" + DeviceName + "/twin/update"

	status := [2]string{strconv.Itoa(int(temperature)) + "C", strconv.Itoa(int(humidity)) + "%"}
	updateMessage := createActualUpdateMessage(status)
	twinUpdateBody, _ := json.Marshal(updateMessage)

	cli.Publish(&client.PublishOptions{
		TopicName: []byte(deviceTwinUpdate),
		QoS:       mqtt.QoS0,
		Message:   twinUpdateBody,
	})
}

//createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(actualValue [2]string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	actualMap := map[string]*MsgTwin{
		DeviceTwinProperties[0]: {Actual: &TwinValue{Value: &actualValue[0]}, Metadata: &TypeMetadata{Type: "Updated temperature"}},
		DeviceTwinProperties[1]: {Actual: &TwinValue{Value: &actualValue[1]}, Metadata: &TypeMetadata{Type: "Updated humidity"}},
	}
	deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}
