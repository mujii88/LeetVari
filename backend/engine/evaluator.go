package engine

import (
	"bytes"
	"context"
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

	cmd:=exec.Command("doceker","exec",workerName,"iverilog",codepath,tbpath,"-o",simpath)

	var compileErr bytes.Buffer
	cmd.Stderr=&compileErr

	if err:=cmd.Run(); err!=nil{
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

	// logs:=simOut.String()





 return Result{Status : "First Handshake", Output: "output", Error: "Not till Now"}

}