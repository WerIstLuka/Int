package main

import (
	"os"
	"io"
	"strings"
)

func HasPipeInput()bool{
	FileInfo, _ := os.Stdin.Stat()
	return FileInfo.Mode() & os.ModeCharDevice == 0
}

func GetArguments() ([]string, []string){
	Options := []string{}
	Numbers := []string{}
	//read input from pipe if it exists
	if HasPipeInput(){
 		bytes, _ := io.ReadAll(os.Stdin)
		pipe := strings.Split((string(bytes)), " ")
		println(pipe)
		for i:=0; i<len(pipe); i++{
			if pipe[i][0:1] == "-"{
				Options = append(Options, pipe[i])
			}else{
				Numbers = append(Numbers, pipe[i])
			}
		}
	}
	//get arguments
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0:1] == "-"{
			Options = append(Options, os.Args[i])
		}else{
			Numbers = append(Numbers, os.Args[i])
		}
	}
	return Options, Numbers
}

func main() {
	GetArguments()
}
