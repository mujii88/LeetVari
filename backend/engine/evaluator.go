package evaluate

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)



type Result struct {
	Status  string `json:"status"`
	Output string `json:"output"`
	Error string `json:"error,omitempty"`
	Metrices Metric `json:"metrics,omitempty"`

}

type Metric struct{
	Number_of_wires int `json:"number_of_wires"`
	Number_of_wire_bits int `json:"number_of_wire_bits"`
	Number_of_public_wires int `json:"number_of_public_wires"`
	Number_of_public_wire_bits int `json:"number_of_public_wire_bits"`
	Number_of_memories int `json:"number_of_memories"`
	Number_of_memory_bits int `json:"number_of_memory_bits"`
	Number_of_processes int `json:"number_of_processes"`
	Number_of_cells int `json:"number_of_cells"`

}


type YosysJSON struct {
	Design struct {
		NumWires       int `json:"num_wires"`
		NumWireBits    int `json:"num_wire_bits"`
		NumPubWires    int `json:"num_pub_wires"`
		NumPubWireBits int `json:"num_pub_wire_bits"`
		NumMemories    int `json:"num_memories"`
		NumMemoryBits  int `json:"num_memory_bits"`
		NumProcesses   int `json:"num_processes"`
		NumCells       int `json:"num_cells"`
	} `json:"design"`
}

func Evaluate(code string , testbench string , workerName string) Result {


	tempDir,err:=os.MkdirTemp("/dev/shm","leetvari_*")
	if err!=nil {
		return Result{Status:"Internal Error", Error:"Failed to create the temporary directory", Output:"null"}
	}
	defer os.RemoveAll(tempDir)


	codepath:=filepath.Join(tempDir,"code.v")
	tbpath:=filepath.Join(tempDir,"tb.v")
	simpath:=filepath.Join(tempDir,"sim.vvp")
	statsPath := filepath.Join(tempDir, "stats.json")


	if err:= os.WriteFile(codepath,[]byte(code),0644); err!=nil{
		return Result{
			Status: "Error in writing the code to temp file",
			Error: err.Error(),
		}
	}


		if err:= os.WriteFile(tbpath,[]byte(testbench),0644); err!=nil{
		return Result{
			Status: "Error in writing the testbench to temp file",
			Error: err.Error(),
		}
	}

	compileCtx, compileCancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer compileCancel()

	cmd:=exec.CommandContext(compileCtx,"docker","exec",workerName,"iverilog",codepath,tbpath,"-o",simpath)

	var compileErr bytes.Buffer
	cmd.Stderr=&compileErr

	if err:=cmd.Run(); err!=nil{
		if compileCtx.Err() == context.DeadlineExceeded {
            return Result{Status: "Compile Error", Error: "Compilation timed out (code too large or complex)"}
        }
		cleanError:=strings.ReplaceAll(compileErr.String(),tempDir+"/","")
		return Result{
			Status: "Compile Error",
			Error:cleanError,
		}
		
	}

	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	cmd2:=exec.CommandContext(ctx,"docker","exec",workerName,"vvp",simpath)

	var simOut bytes.Buffer
	cmd2.Stdout=&simOut

	if err:=cmd2.Run(); err!=nil{

		if ctx.Err()==context.DeadlineExceeded{
			return Result{
				Status: "TIme Limit Exceed",

			}
		}
		return Result{
			Status: "Run time error",
			Error:err.Error(),
		}

	}

	logs:=simOut.String()

	if strings.Contains(logs,"LEETVARI_PASS") {
		cleanoutput:=strings.ReplaceAll(logs,tempDir+"/","")
		yosys_ctx,yosys_cancel:=context.WithTimeout(context.Background(),5*time.Second)
		defer yosys_cancel()


		
        yosysScript := "synth -auto-top; tee -o " + statsPath + " stat -json"
        
        cmd3 := exec.CommandContext(yosys_ctx, "docker", "exec", workerName, "yosys", "-q", "-p", yosysScript, codepath)
        
        if err := cmd3.Run(); err != nil {
            return Result{Status: "INTERNAL_ERROR", Error: err.Error(), Output:"null"}
        }

        statsBytes, err := os.ReadFile(statsPath)
        if err != nil {
            return Result{Status: "INTERNAL_ERROR", Error: "Failed to read stats.json file from /dev/shm"}
        }

        var parsedYosys YosysJSON
        if err := json.Unmarshal(statsBytes, &parsedYosys); err != nil {
            return Result{
                Status: "INTERNAL_ERROR", 
                Error:  err.Error() + "\nRAW YOSYS OUTPUT:\n" + string(statsBytes),
            }
        }

        finalMetrics := Metric{
            Number_of_wires:            parsedYosys.Design.NumWires,
            Number_of_wire_bits:        parsedYosys.Design.NumWireBits,
            Number_of_public_wires:     parsedYosys.Design.NumPubWires,
            Number_of_public_wire_bits: parsedYosys.Design.NumPubWireBits,
            Number_of_memories:         parsedYosys.Design.NumMemories,
            Number_of_memory_bits:      parsedYosys.Design.NumMemoryBits,
            Number_of_processes:        parsedYosys.Design.NumProcesses,
            Number_of_cells:            parsedYosys.Design.NumCells,
        }

		return Result{
			Status:   "Accepted",
			Output:   cleanoutput,
			Metrices: finalMetrics,
		}

		
	}

		cleanOutput := strings.ReplaceAll(logs, tempDir+"/", "")
		return Result{
			Status: "WRONG_ANSWER",
			Output: cleanOutput,
	    }


}