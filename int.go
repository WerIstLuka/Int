package main

import (
	"os"
	"io"
	"strings"
	"slices"
	"regexp"
	"strconv"
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

func GetInt(char string)uint64{
	switch char{
		case "b":
			return 2
		case "o":
			return 8
		case "x":
			return 16
		default:
			return 0
	}
	return 0
}


func Parser(Options []string)(uint64, uint64){
	digitCheck := regexp.MustCompile(`^[0-9]+$`)
	var GotInputBase bool = false
	var InputBase uint64 = 0
	var OutputBase uint64 = 0
	for i:=0; i<len(Options); i++{
		Option := Options[i]
		if len(Option) < 2{
			println("Error: invalid option \"", Option, "\"")
		}
		if InputBase != 0 && slices.Contains([]string{"-B", "-b", "-o", "-x"}, Option){
			println("Error: Input base was given twice")
			os.Exit(1)
		}
		if Option[1:2] == "B" && digitCheck.MatchString(Option[2:len(Option)-1]){
			InputBase, _ = strconv.ParseUint(Option[2:len(Option)-1], 10, 64)
			GotInputBase = true
		}
		if len(Option) == 2 && slices.Contains([]string{"b", "o", "x"}, Option[1:2]){
			InputBase = GetInt(Option[1:2])
		}
		if GotInputBase && InputBase == 0{
			println("Error: input base can't be 0")
		}
	println(InputBase)
	println(OutputBase)
	}
	return InputBase, OutputBase
}


func main() {
	Options, _ := GetArguments()
	Parser(Options)
}
