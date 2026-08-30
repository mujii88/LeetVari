package engine



type Result struct {
	Status  string `json:"status"`
	Output string `json:"output"`
	Error string `json:"error,omitempty"`

}

func Evaluate(code string , testbench string , workerName string) Result {


 return Result{Status : "First Handshake", Output: "output", Error: "Not till Now"}

}