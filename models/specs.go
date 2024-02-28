package models

type MouseSpecs struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
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
	DPI          string `bson:"dpi"`
}
type KeyBoardSpecs struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
	Form_Factor string `bson:"form_factor"`
	PCB         string `bson:"PCB"`
	Height      string `bson:"height"`
	Length      string `bson:"length"`
	Switches    string `bson:"Switches"`
	RGB         string `bson:"RGB"`
	Width       string `bson:"width"`
	Weight      string `bson:"weight"`
}

type HeadsetSpecs struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
	Headset_Type     string `bson:"headset_type"`
	Cable_Length     string `bson:"cable_length"`
	Microphone       string `bson:"microphone"`
	Connection       string `bson:"connection"`
	Noise_Cancelling string `bson:"noise_cancelling"`
	Weight           string `bson:"weight"`
}

type GPU struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
	NVIDIA_CUDA_Cores string `bson:"nvidia_cuda_cores"`
	Memory_Size       string `bson:"memory_size"`
	Boost_Clock       string `bson:"boost_clock"`
	Memory_Type       string `bson:"memory_type"`
}

type MousePad struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
	Height         string `bson:"height"`
	Thickness      string `bson:"thickness"`
	Material       string `bson:"material"`
	Length         string `bson:"length"`
	Stitched_edges string `bson:"stitched_edges"`
	Glide          string `bson:"glide"`
}

type CPU struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
	Socket           string `bson:"socket"`
	Threads          string `bson:"threads"`
	Core_Speed_Base  string `bson:"core_speed_base"`
	Cores            string `bson:"cores"`
	TDP              string `bson:"TDP"`
	Core_Speed_Boost string `bson:"core_speed_boost"`
}
type Monitor struct {
	Type string 	`bson:"type"`
	PID       string `bson:"pID"`
	Size         string `bson:"size"`
	Aspect_Ratio string `bson:"aspect_ratio"`
	G_Sync       string `bson:"g_sync"`
	Panel_Tech   string `bson:"panel_tech"`
	Resolution   string `bson:"resolution"`
	Refresh_Rate string `bson:"refresh_rate"`
	FreeSync     string `bson:"free_sync"`
}
