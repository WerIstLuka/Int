package main

import (
	"os"
	"io"
	"strings"
	"slices"
	"regexp"
	"strconv"
)

var Version string = "2.0-dev"

func Help(){
	println(`Convert any base to any other
Usage:
	int [Options] [Integers]
Options:
	-b	add binary prefix to all integers
	-x	add hexadecimal prefix to all integers
	-o	add octal prefix to all integers
	-Bx	where x is the base of the integers
	-Ox	set the output to any base, also works with b, x and o
	-h	show this Help text`)
	os.Exit(0)
}

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
}

func GetBase(Option string)uint64{
	var Base uint64
	digitCheck := regexp.MustCompile(`^[0-9]+$`)
	if digitCheck.MatchString(Option){
		Base, _ = strconv.ParseUint(Option, 10, 64)
	} else {
		Base = GetInt(Option)
	}
	if Base == 0{
		println("Error: invalid base:", Option)
		os.Exit(1)
	}
	return Base
}


func Parser(Options []string)(uint64, uint64){
	var GotInputBase bool = false
	var GotOutputBase bool = false
	var InputBase uint64 = 0
	var OutputBase uint64 = 0
	for i:=0; i<len(Options); i++{
		Option := Options[i]
		println(Option)
		if len(Option) < 2{
			println("Error: invalid option:", Option)
			os.Exit(1)
		}
		if InputBase != 0 && slices.Contains([]string{"-B", "-b", "-o", "-x"}, Option){
			println("Error: Input base was given twice")
			os.Exit(1)
		} else if OutputBase != 0 && Option[2:len(Option)] == "O"{
			println("Error: Output base was given twice")
			os.Exit(1)
		} else if len(Option) == 2 && slices.Contains([]string{"b", "o", "x"}, Option[1:2]){
			InputBase = GetInt(Option[1:2])
		} else if GotInputBase && InputBase == 0{
			println("Error: Input base can't be 0")
			os.Exit(1)
		} else if GotOutputBase && OutputBase ==0{
			println("Error: Output base can't be 0")
			os.Exit(1)
		} else if Option == "-h" || Option == "--help"{
			Help()
		} else if Option == "-v" || Option == "--version"{
			println(Version)
			os.Exit(0)
		} else if slices.Contains([]string{"B", "O"}, Option[1:2]) == false{
			println("Error: invalid option:", Option)
			os.Exit(1)
		} else if len(Option) < 3{
			println("Error: too few arguments:", Option)
			println(Option[1:2])
			os.Exit(1)
		} else if Option[1:2] == "B"{
			InputBase = GetBase(Option[2:len(Option)])
			GotInputBase = true
		} else if Option[1:2] == "O"{
			OutputBase = GetBase(Option[2:len(Option)])
			GotOutputBase = true
		}
	}
	println("In:", InputBase)
	println("Out:", OutputBase)
	return InputBase, OutputBase
}

func ConvertNumbers(Numbers []string, InputBase uint64) {
	ConvertedNumbers := []uint64{}
	for i:=0; i<len(Numbers); i++{
		NumberStr := Numbers[i]
		println(NumberStr)
		if len(NumberStr) >= 2{
			switch NumberStr[0:2]{
				case "0b":
					InputBase = 2
			}
		}
		NumberInt, err := strconv.ParseUint(NumberStr, 10, 64)
		println(err)
		if err == nil{
			println("Operation not possible:", NumberStr)
			os.Exit(1)
		}
		println(NumberInt)
		ConvertedNumbers = append(ConvertedNumbers, NumberInt)
	}
}

func main() {
	Options, Numbers := GetArguments()
	InputBase, _ := Parser(Options)
	if len(Numbers) == 0{
		println("no input")
		os.Exit(0)
	}
	ConvertNumbers(Numbers, InputBase)
}
