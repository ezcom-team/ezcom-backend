package models

type MouseSpecs struct {
	ProductID    string `bson:"productId"`
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
