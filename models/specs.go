package models

type MouseSpecs struct {
	SpecsID      string `bson:"specId"`
	Sensor       string `bson:"sensor"`
	ButtonSwitch string `bson:"buttonSwitch"`
	Connection   string `bson:"connection"`
	Length       string `bson:"legth"`
	Weight       string `bson:"weight"`
	PollingRate  string `bson:"pollingRate"`
	ButtonForce  string `bson:"buttonForce"`
	Shape        string `bson:"shape"`
	Height       string `bson:"height"`
	Width        string `bson:"width"`
}
type KeyBoardSpecs struct {
	SpecsID     string `bson:"specId"`
	Form_Factor string `bson:"form_factor"`
	PCB         string `bson:"PCB"`
	Height      string `bson:"height"`
	Length      string `bson:"legth"`
	Switches    string `bson:"Switches"`
	RGB         string `bson:"RGB"`
	Width       string `bson:"width"`
	Weight      string `bson:"weight"`
}

type HeadsetSpecs struct {
	SpecsID          string `bson:"specId"`
	Headset_Type     string `bson:"headset_type"`
	Cable_Length     string `bson:"cable_length"`
	Microphone       string `bson:"microphone"`
	Connection       string `bson:"connection"`
	Noise_Cancelling string `bson:"noise_cancelling"`
	Weight           string `bson:"weight"`
}
