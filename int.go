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

var version string = "2.2.0"

func help(){
	fmt.Println(`Convert any base to any other
Usage:
	int [Options] [Integers]
Options:
	-b  add binary prefix to all integers
	-x  add hexadecimal prefix to all integers
	-o  add octal prefix to all integers
	-Bx where x is the base of the integers
	-Ox set the output to any base, also works with b, x and o
 	-l  use the same characters for bases below and above 36
 	-ax aligns output, X is an integer and aligns to the next highest multiple
 	-s  split output at multiples of the alignment, requires -ax
	-h  show this Help text`)
	os.Exit(0)
}

func getArguments() ([]string, []string){
	options := []string{}
	numbers := []string{}
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
					options = append(options, strings.TrimSpace(line[j]))
				}else{
					numbers = append(numbers, strings.TrimSpace(line[j]))
				}
			}
		}
	}
	//get arguments
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] == '-'{
			options = append(options, os.Args[i])
		}else{
			numbers = append(numbers, os.Args[i])
		}
	}
	if len(options) == 0 && len(numbers) == 0{
		help()
	}
	return options, numbers
}

func getInt(char byte)int64{
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

func getBase(option string)int64{
	var base int64
	var isValid bool = true
	for i:=0;i<len(option);i++{
		if !strings.Contains("0123456789", string(option[i])){
			isValid = false
		}
	}
	if isValid{
		base, _ = strconv.ParseInt(option, 10, 64)
	} else {
		if len(option) != 1{
			base = 0
		} else {
			base = getInt(option[0])
		}
	}
	if base == 0 || base > 62{
		llg.PrintError("Invalid base:", option)
	}
	return base
}

func checkAlign(opt string, align *int, gotInt *bool) bool {
	var alignment int64 = 0
	var isAligned bool = false
	if opt[:2] == "-a" {
		isAligned = true
		alignment, _ = strconv.ParseInt(opt[2:], 0, 32)
	} else if len(opt) >= 7 && opt[:7] == "--align" {
		isAligned = true
		alignment, _ = strconv.ParseInt(opt[7:], 0, 32)
	}
	if !isAligned { return false }
	if *align != 0 {
		llg.PrintError("Align was given twice")
	}
	if alignment == 0 {
		*align = 1
	} else {
		*align = int(alignment)
		
	}
	return true
}

func parser(options []string)(int, int64, bool, int, bool){
	var gotInputBase bool = false
	var gotOutputBase bool = false
	var inputBase int64 = 0
	var outputBase int64 = 0
	var forceLong bool = false
	var align int = 0
	var split bool = false
	var splitGotInt bool = false
	for i:=0; i<len(options); i++{
		option := options[i]
		if len(option) < 2{
			llg.PrintError("Invalid option:", option)
		}
		if checkAlign(option, &align, &splitGotInt) {
			continue
		} else if gotInputBase && slices.Contains([]byte{'B', 'b', 'o', 'x'}, option[1]){
			llg.PrintError("Input base was given twice")
		} else if gotOutputBase && option[1] == 'O'{
			llg.PrintError("Output base was given twice")
		} else if len(option) == 2 && slices.Contains([]byte{'b', 'o', 'x'}, option[1]){
			inputBase = getInt(option[1])
			gotInputBase = true
		} else if option == "-h" || option == "--help"{
			help()
		} else if option == "-v" || option == "--version"{
			fmt.Println(version)
			os.Exit(0)
		} else if option == "-l" || option == "--long"{
			forceLong = true
		} else if option == "-s" || option == "--split"{
			split = true
		} else if slices.Contains([]byte{'B', 'O'}, option[1]) == false{
			llg.PrintError("Invalid option:", option)
		} else if len(option) < 3{
			llg.PrintError("Too few arguments:", option)
		} else if option[1] == 'B'{
			inputBase = getBase(option[2:len(option)])
			gotInputBase = true
		} else if option[1] == 'O'{
			outputBase = getBase(option[2:len(option)])
			gotOutputBase = true
		}
	}
	if gotInputBase && inputBase <= 1{
		llg.PrintError("Input base can't be 1 or less")
	} else if gotOutputBase && outputBase <= 1{
		llg.PrintError("Output base can't be 1 or less")
	}
	if split {
		if splitGotInt {
			llg.PrintError("Split needs align with an integer")
		}
	}
	return int(inputBase), outputBase, forceLong, align, split
}

func errorOperationNotPossible(operation string){
	llg.PrintError("Operation not possible:", operation)
}

func convertNumbers(num string, inputBase int, outputBase int64, forceLong bool) string{
	if inputBase == 0 && len(num) >= 2{
		inputBase = int(getInt(num[1]))
		if inputBase != 0{
			if len(num) == 2{
				llg.PrintError("No number given:", num)
			}
			num = num[2:len(num)]
		}
	}
	if inputBase == 0{
		inputBase = 10
	}
	if outputBase == 0{
		outputBase = 10
	}

	var digitsShort string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var digitsLong string  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var digitList []string

	if inputBase <= 35{
		digitList = strings.Split(digitsShort, "")
	} else {
		digitList = strings.Split(digitsLong, "")
	}
	var index int
	for i:=0;i<len(num);i++{
		index = slices.Index(digitList, string(num[i]))
		if index == -1 || forceLong{
			index = slices.Index(strings.Split(digitsLong, ""), string(num[i]))
			if index == -1{
				llg.PrintError("Illegal character:", string(num[i]))
			}
		}
		if index >= inputBase{
			llg.PrintError("Number", num, "is not valid for base", inputBase)
		}
	}
	var outputDigits *string
	if outputBase <= 35 && !forceLong{
		outputDigits = &digitsShort
	} else {
		outputDigits = &digitsLong
	}

	bigNum := new(big.Int)
	bigNum.SetString(num, inputBase)

	bigOutputBase := big.NewInt(outputBase)
	var output = []string{}
	indexBig := new(big.Int)
	if bigNum.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}
	for bigNum.Cmp(big.NewInt(0)) != 0 {
		output = append(output, string(string(*outputDigits)[indexBig.Mod(bigNum, bigOutputBase).Int64()]))
		bigNum.Div(bigNum, bigOutputBase)
	}
	slices.Reverse(output)
	return strings.Join(output, "")
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
	options, numbers := getArguments()
	inputBase, outputBase, forceLong, align, split := parser(options)
	if len(numbers) == 0{
		fmt.Println("no input")
		os.Exit(0)
	}
	convertedNumbers := []string{}
	for i:=0;i<len(numbers);i++{
		convertedNumbers = append(convertedNumbers, convertNumbers(numbers[i], inputBase, outputBase, forceLong))
	}
	if align == 0 {
		for i:=0;i<len(convertedNumbers);i++{
			fmt.Println(convertedNumbers[i])
		}
		os.Exit(0)
	}
	printAligned(&align, &convertedNumbers, &split)
}
