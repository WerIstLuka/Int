#!/bin/python3

#developed at https://github.com/WerIstLuka/int

import sys
import numpy

Version = "1.5.4"

def Help():
	print("""\
Convert any base to any other
Usage:
	int [options] [integers]
options:
	-b		add binary prefix to all integers
	-x		add hexadecimal prefix to all integers
	-o		add octal prefix to all integers
	-Bx		where x is the base of the integers
	-Ox		set the output to any base, also works with b, x and o
	-h		show this Help text\
""")
	quit()

def Exit_Too_Many_Prefixes():
	print("Conflicting Prefixes")
	quit()

def GetInt(Arguments, i):
	match Arguments[i][2:]:
		case "x":
			OutputBase = 16
		case "b":
			OutputBase = 2
		case "o":
			OutputBase = 8
		case _:
			try:
				OutputBase = int(Arguments[i][2:])
			except Exception:
				print(f"Not a valid base: {Arguments[i]}")
				quit()
	return OutputBase

def CheckOptions(Arguments):
	InputBase = 10
	OutputBase = 10
	OutputArgument = ""
	InputArgument = ""
	for i in range(len(Arguments)):
		match Arguments[i][:2]:
			case "-b":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				InputBase = 2
				InputArgument = "-b"
			case "-x":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				InputBase = 16
				InputArgument = "-x"
			case "-o":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				InputBase = 8
				InputArgument = "-o"
			case "-B":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				if Arguments[i][2:] == "":
					print(f"No Base given: {Arguments[i]}")
					quit()
				InputBase = GetInt(Arguments, i)
				InputArgument = Arguments[i]
			case "-O":
				if OutputBase != 10:
					print("Output Base given twice")
					quit()
				if Arguments[i][2:] == "":
					print(f"No Base given: {Arguments[i]}")
					quit()
				OutputBase = GetInt(Arguments, i)
				OutputArgument = Arguments[i]
			case _:
				pass
	if "-h" in Arguments or "--help" in Arguments:
		Help()
	if "--version" in Arguments:
		print(f"int {Version}")
		quit()
	return(InputBase, OutputBase, OutputArgument, InputArgument)

def GetIntegers(Arguments, InputBase, OutputBase, OutputArgument, InputArgument):
	if OutputBase != 10:
		Arguments.remove(OutputArgument)
	DecimalIntegers = []
	if InputBase != 10:
		Arguments.remove(InputArgument)
		for i in range(len(Arguments)):
			if Arguments[i][0] == "-":
				print(f"unknown option: {Arguments[i]}")
				quit()
			try:
				int(Arguments[i], InputBase)
			except Exception:
				print(f"Operation not possible: {Arguments[i]}")
				quit()
			DecimalIntegers.append(int(Arguments[i], InputBase))
	else:
		for i in range(len(Arguments)):
			if Arguments[i][0] == "-":
				print(f"unknown option: {Arguments[i]}")
				quit()
			match Arguments[i][:2]:
				case "0b":
					DecimalIntegers.append(int(Arguments[i][2:], 2))
				case "0x":
					DecimalIntegers.append(int(Arguments[i][2:], 16))
				case _:
					try:
						int(Arguments[i])
					except Exception:
						print(f"Operation not possible: {Arguments[i]}")
						quit()
					DecimalIntegers.append(int(Arguments[i]))
	if len(Arguments) == 0:
		print("no input given")
		quit()
	return DecimalIntegers

def OutputIntegers(DecimalIntegers, OutputBase):
	for i in range(len(DecimalIntegers)):
		print(numpy.base_repr(DecimalIntegers[i], int(OutputBase)))

# get pipes before arguments (read left to right)
Arguments = []
Debugger = False
if not Debugger: #debugger hangs when trying to read from stdin
	if not sys.stdin.isatty():
		for line in sys.stdin:
			line = line.rstrip()
			Arguments += line.split(" ")
for i in range(len(sys.argv[1:])):
	Arguments.append(sys.argv[1:][i])

if len(Arguments) == 0:
	Help()

Options = CheckOptions(Arguments)
DecimalOutput =  GetIntegers(Arguments, Options[0], Options[1], Options[2], Options[3])
OutputIntegers(DecimalOutput, Options[1])
