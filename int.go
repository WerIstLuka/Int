package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"slices"
	"strconv"
	"math/big"
	"github.com/WerIstLuka/llg"
	"golang.org/x/term"
)

var Version string = "2.2.0"

func Help(){
	fmt.Println(`Convert any base to any other
Usage:
	int [Options] [Integers]
Options:
	-b	add binary prefix to all integers
	-x	add hexadecimal prefix to all integers
	-o	add octal prefix to all integers
	-Bx	where x is the base of the integers
	-Ox	set the output to any base, also works with b, x and o
 	-l	use the same characters for bases below and above 36
 	-ax aligns output, X is an integer and align to the next highest multiple
 	-s  split output at multiples of the alignment, requires -ax
	-h	show this Help text`)
	os.Exit(0)
}

func GetArguments() ([]string, []string){
	Options := []string{}
	Numbers := []string{}
	//read input from pipe if it exists
	if !term.IsTerminal(0) {
 		bytes, _ := io.ReadAll(os.Stdin)
		lines := strings.Split((string(bytes)), "\n")
		for i:=0; i<len(lines); i++{
			line := strings.Split(lines[i], " ")
			for j:=0; j<len(line); j++{
				if len(line[j]) == 0{
					continue
				}
				if line[j][0] == '-'{
					Options = append(Options, strings.TrimSpace(line[j]))
				}else{
					Numbers = append(Numbers, strings.TrimSpace(line[j]))
				}
			}
		}
	}
	//get arguments
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] == '-'{
			Options = append(Options, os.Args[i])
		}else{
			Numbers = append(Numbers, os.Args[i])
		}
	}
	if len(Options) == 0 && len(Numbers) == 0{
		Help()
	}
	return Options, Numbers
}

func GetInt(char byte)int64{
	switch char{
		case 'b':
			return 2
		case 'o':
			return 8
		case 'x':
			return 16
		default:
			return 0
	}
}

func GetBase(Option string)int64{
	var Base int64
	var IsValid bool = true
	for i:=0;i<len(Option);i++{
		if !strings.Contains("0123456789", string(Option[i])){
			IsValid = false
		}
	}
	if IsValid{
		Base, _ = strconv.ParseInt(Option, 10, 64)
	} else {
		if len(Option) != 1{
			Base = 0
		} else {
			Base = GetInt(Option[0])
		}
	}
	if Base == 0 || Base > 62{
		llg.PrintError("Invalid base:", Option)
	}
	return Base
}

func checkAlign(opt string, align *int) bool {
	var Alignment int64 = 0
	var isAligned bool = false
	if opt[:2] == "-a" {
		isAligned = true
		Alignment, _ = strconv.ParseInt(opt[2:], 0, 32)
	} else if len(opt) >= 7 && opt[:7] == "--align" {
		isAligned = true
		Alignment, _ = strconv.ParseInt(opt[7:], 0, 32)
	}
	if !isAligned { return false }
	if *align != 0 {
		llg.PrintError("Align was given twice")
	}
	if Alignment == 0 {
		*align = 1
	} else {
		*align = int(Alignment)
	}
	return true
}

func Parser(Options []string)(int, int64, bool, int, bool){
	var GotInputBase bool = false
	var GotOutputBase bool = false
	var InputBase int64 = 0
	var OutputBase int64 = 0
	var ForceLong bool = false
	var Align int = 0
	var split bool = false
	for i:=0; i<len(Options); i++{
		Option := Options[i]
		if len(Option) < 2{
			llg.PrintError("Invalid option:", Option)
		}
		if checkAlign(Option, &Align) {
			continue
		} else if GotInputBase && slices.Contains([]byte{'B', 'b', 'o', 'x'}, Option[1]){
			llg.PrintError("Input base was given twice")
		} else if GotOutputBase && Option[1] == 'O'{
			llg.PrintError("Output base was given twice")
		} else if len(Option) == 2 && slices.Contains([]byte{'b', 'o', 'x'}, Option[1]){
			InputBase = GetInt(Option[1])
			GotInputBase = true
		} else if Option == "-h" || Option == "--help"{
			Help()
		} else if Option == "-v" || Option == "--version"{
			fmt.Println(Version)
			os.Exit(0)
		} else if Option == "-l" || Option == "--long"{
			ForceLong = true
		} else if Option == "-s" || Option == "--split"{
			split = true
		} else if slices.Contains([]byte{'B', 'O'}, Option[1]) == false{
			llg.PrintError("Invalid option:", Option)
		} else if len(Option) < 3{
			llg.PrintError("Too few arguments:", Option)
		} else if Option[1] == 'B'{
			InputBase = GetBase(Option[2:len(Option)])
			GotInputBase = true
		} else if Option[1] == 'O'{
			OutputBase = GetBase(Option[2:len(Option)])
			GotOutputBase = true
		}
	}
	if GotInputBase && InputBase <= 1{
		llg.PrintError("Input base can't be 1 or less")
	} else if GotOutputBase && OutputBase <= 1{
		llg.PrintError("Output base can't be 1 or less")
	}
	if split {
		if Align == 0 || Align == 1 {
			llg.PrintError("Split needs align with an integer")
		}
	}
	return int(InputBase), OutputBase, ForceLong, Align, split
}

func ErrorOperationNotPossible(Operation string){
	llg.PrintError("Operation not possible:", Operation)
}

func ConvertNumbers(Num string, InputBase int, OutputBase int64, ForceLong bool) string{
	if InputBase == 0 && len(Num) >= 2{
		InputBase = int(GetInt(Num[1]))
		if InputBase != 0{
			if len(Num) == 2{
				llg.PrintError("No number given:", Num)
			}
			Num = Num[2:len(Num)]
		}
	}
	if InputBase == 0{
		InputBase = 10
	}
	if OutputBase == 0{
		OutputBase = 10
	}

	var DigitsShort string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var DigitsLong string  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var DigitList []string

	if InputBase <= 35{
		DigitList = strings.Split(DigitsShort, "")
	} else {
		DigitList = strings.Split(DigitsLong, "")
	}
	var Index int
	for i:=0;i<len(Num);i++{
		Index = slices.Index(DigitList, string(Num[i]))
		if Index == -1 || ForceLong{
			Index = slices.Index(strings.Split(DigitsLong, ""), string(Num[i]))
			if Index == -1{
				llg.PrintError("Illegal character:", string(Num[i]))
			}
		}
		if Index >= InputBase{
			llg.PrintError("Number", Num, "is not valid for base", InputBase)
		}
	}
	var OutputDigits *string
	if OutputBase <= 35 && !ForceLong{
		OutputDigits = &DigitsShort
	} else {
		OutputDigits = &DigitsLong
	}

	BigNum := new(big.Int)
	BigNum.SetString(Num, InputBase)

	BigOutputBase := big.NewInt(OutputBase)
	var Output = []string{}
	IndexBig := new(big.Int)
	if BigNum.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}
	for BigNum.Cmp(big.NewInt(0)) != 0 {
		Output = append(Output, string(string(*OutputDigits)[IndexBig.Mod(BigNum, BigOutputBase).Int64()]))
		BigNum.Div(BigNum, BigOutputBase)
	}
	slices.Reverse(Output)
	return strings.Join(Output, "")
}

func printSplit(num *string, size int, align *int) {
	var index int = 0
	var chunk int = *align
	for index < size {
		if index + chunk >= size {
			chunk = size - index
		}
		fmt.Print((*num)[index:index + chunk])
		index += chunk
		if index != size {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

func printAligned(align *int, output *[]string, split *bool) {
	var longest = llg.MaxStringLength(*output)
	var outLength int = llg.NextMultipleInt(longest, *align)
	var num string
	var curSize int
	for i:=0; i<len(*output); i++ {
		num = (*output)[i]
		curSize = len(num)
		num = strings.Repeat("0", outLength-curSize) + num
		if !*split {
			fmt.Println(num)
		} else {
			printSplit(&num, len(num), align)
		}
	}
}

func main() {
	Options, Numbers := GetArguments()
	InputBase, OutputBase, ForceLong, Align, split := Parser(Options)
	if len(Numbers) == 0{
		fmt.Println("no input")
		os.Exit(0)
	}
	ConvertedNumbers := []string{}
	for i:=0;i<len(Numbers);i++{
		ConvertedNumbers = append(ConvertedNumbers, ConvertNumbers(Numbers[i], InputBase, OutputBase, ForceLong))
	}
	if Align == 0 {
		for i:=0;i<len(ConvertedNumbers);i++{
			fmt.Println(ConvertedNumbers[i])
		}
		os.Exit(0)
	}
	printAligned(&Align, &ConvertedNumbers, &split)
}
