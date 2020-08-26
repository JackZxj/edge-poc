package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/stianeikeland/go-rpio"
)

const (
	DeviceETPrefix            = "$hw/events/device/"
	DeviceETStateUpdateSuffix = "/state/update"
	TwinETUpdateSuffix        = "/twin/update"
	TwinETCloudSyncSuffix     = "/twin/cloud_updated"
	TwinETGetResultSuffix     = "/twin/get/result"
	TwinETGetSuffix           = "/twin/get"
)

var rgb []string
var defaultRGBValue []int
var rgbPinNumber []int
var Token_client MQTT.Token
var ClientOpts *MQTT.ClientOptions
var Client MQTT.Client
var wg sync.WaitGroup
var deviceTwinResult DeviceTwinUpdate
var deviceID string
var modelName string
var DeviceName string
var MQTTURL string

//usage is responsible for setting up the default settings of all defined command-line flags for glog.
func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

//init for getting command line arguments for glog and initiating the MQTT connection
func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()

	modelName = "RGB-LIGHT"
	DeviceName = "rgb-light-device"
	MQTTURL = "tcp://127.0.0.1:1884"
	rgb = []string{"red-pwm", "green-pwm", "green-pwm"}
	defaultRGBValue = []int{50, 50, 50}
	rpin, _ := strconv.Atoi(os.Getenv("RPIN"))
	gpin, _ := strconv.Atoi(os.Getenv("RPIN"))
	bpin, _ := strconv.Atoi(os.Getenv("RPIN"))
	rgbPinNumber = []int{rpin, gpin, bpin}

	glog.Info("Init MQTT client...")
	ClientOpts = HubClientInit(MQTTURL, "eventbus", "", "")
	Client = MQTT.NewClient(ClientOpts)
	if Token_client = Client.Connect(); Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.Connect() Error is ", Token_client.Error())
	}
	err := LoadConfigMap()
	if err != nil {
		glog.Error(errors.New("Error while reading from config map " + err.Error()))
		os.Exit(1)
	}
}

func main() {
	changeDeviceState("online")
	updateMessage := createActualUpdateMessage(defaultRGBValue)
	for {
		equateTwinValue(updateMessage)
	}
}

// HubclientInit create mqtt client config
func HubClientInit(server, clientID, username, password string) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientID).SetCleanSession(true)
	if username != "" {
		opts.SetUsername(username)
		if password != "" {
			opts.SetPassword(password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)
	return opts
}

// LoadConfigMap load config from /opt/kubeedge/deviceProfile.json to get deviceID
func LoadConfigMap() error {
	readConfigMap := DeviceProfile{}
	err := readConfigMap.ReadFromConfigMap()
	if err != nil {
		return errors.New("Error while reading from config map " + err.Error())
	}
	for _, device := range readConfigMap.DeviceInstances {
		if strings.ToUpper(device.Model) == modelName && strings.ToUpper(device.Name) == strings.ToUpper(DeviceName) {
			deviceID = device.ID
		}
	}
	return nil
}

//changeDeviceState function is used to change the state of the device
func changeDeviceState(state string) {
	glog.Info("Changing the state of the device to" + state)
	var deviceStateUpdateMessage DeviceStateUpdate
	deviceStateUpdateMessage.State = state
	stateUpdateBody, err := json.Marshal(deviceStateUpdateMessage)
	if err != nil {
		glog.Error("Error:   ", err)
	}
	deviceStatusUpdate := DeviceETPrefix + deviceID + DeviceETStateUpdateSuffix
	Token_client = Client.Publish(deviceStatusUpdate, 0, false, stateUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.publish() Error in device state update  is ", Token_client.Error())
	}
}

//createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(slice []int) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	// actualMap := map[string]*MsgTwin{
	// 	"red-pwm":   {Actual: &TwinValue{Value: redPwm}, Metadata: &TypeMetadata{Type: "Updated"}},
	// 	"green-pwm": {Actual: &TwinValue{Value: greenPwm}, Metadata: &TypeMetadata{Type: "Updated"}},
	// 	"blue-pwm":  {Actual: &TwinValue{Value: bluePwm}, Metadata: &TypeMetadata{Type: "Updated"}},
	// }
	actualMap := map[string]*MsgTwin{}
	for i, color := range rgb {
		v := MsgTwin{Actual: &TwinValue{Value: slice[i]}, Metadata: &TypeMetadata{Type: "Updated"}}
		actualMap[color] = &v
	}

	deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}

//changeTwinValue sends the updated twin value to the edge through the MQTT broker
func changeTwinValue(updateMessage DeviceTwinUpdate) {
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		glog.Error("Error:   ", err)
	}
	deviceTwinUpdate := DeviceETPrefix + deviceID + TwinETUpdateSuffix
	Token_client = Client.Publish(deviceTwinUpdate, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.publish() Error in device twin update is ", Token_client.Error())
	}
}

//getTwin function is used to get the device twin details from the edge
func getTwin(updateMessage DeviceTwinUpdate) {
	getTwin := DeviceETPrefix + deviceID + TwinETGetSuffix
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		glog.Error("Error:   ", err)
	}
	Token_client = Client.Publish(getTwin, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.publish() Error in device twin get  is ", Token_client.Error())
	}
}

// OnSubMessageReceived callback function which is called when message is received
func OnSubMessageReceived(client MQTT.Client, message MQTT.Message) {
	err := json.Unmarshal(message.Payload(), &deviceTwinResult)
	if err != nil {
		glog.Error("Error in unmarshalling:  ", err)
	}
}

//subscribe function subscribes  the device twin information through the MQTT broker
func subscribe() {
	for {
		getTwinResult := DeviceETPrefix + deviceID + TwinETGetResultSuffix
		Token_client = Client.Subscribe(getTwinResult, 0, OnSubMessageReceived)
		if Token_client.Wait() && Token_client.Error() != nil {
			glog.Error("subscribe() Error in device twin result get  is ", Token_client.Error())
		}
		time.Sleep(1 * time.Second)
		if deviceTwinResult.Twin != nil {
			wg.Done()
			break
		}
	}
}

// isUpdated Check whether the actual value is synchronized to the expectation
func isUpdated() bool {
	updated := true
	for _, color := range rgb {
		if deviceTwinResult.Twin[color].Expected != nil &&
			((deviceTwinResult.Twin[color].Actual == nil && deviceTwinResult.Twin[color].Expected != nil) ||
				compareValue(color)) {
			updated = false
			break
		}
	}
	return updated
}

// compareValue compare Expected and Actual values, return false when Expected.Value == Expected.Value
func compareValue(key string) bool {
	ev, err1 := (*deviceTwinResult.Twin[key]).Expected.Value.(float64)
	av, err2 := (*deviceTwinResult.Twin[key]).Expected.Value.(float64)
	if err1 == true && err2 == true {
		if ev == av {
			return false
		}
	}
	return true
}

//equateTwinValue is responsible for equating the actual state of the device to the expected state that has been set
func equateTwinValue(updateMessage DeviceTwinUpdate) {
	glog.Info("Watching on the device twin values for device: ", DeviceName)
	wg.Add(1)
	go subscribe()
	getTwin(updateMessage)
	wg.Wait()
	if !isUpdated() {
		msg, err := json.Marshal(deviceTwinResult)
		if err != nil {
			glog.Error("deviceTwinResult Marshal Error:", err)
		}
		glog.Info("DeviceTwinResult : ", msg)
		glog.Info("Equating the actual  value to expected value")
		pwmValues := make([]int, len(rgb))
		for i, color := range rgb {
			expectedValue := (*deviceTwinResult.Twin[color]).Expected.Value.(float64)
			setPWM(rgbPinNumber[i], uint32(expectedValue))
			pwmValues = append(pwmValues, int(expectedValue))
		}
		updateMessage = createActualUpdateMessage(pwmValues)
		changeTwinValue(updateMessage)
	} else {
		glog.Info("Actual values are in sync with Expected value")
	}
}

// setPWM update light intensity
func setPWM(pinNumber int, pwm uint32) {
	pin := rpio.Pin(pinNumber)
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	pin.Mode(rpio.Pwm)
	pin.Freq(50)
	pin.DutyCycle(pwm, 100)
}
