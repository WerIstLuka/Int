# IMPORANT
this is currently work in progress

i decided to rewrite this in Go because python is slow

the python branch will be discontinued after the go version has feature parity

as of now the build script doesnt work

it will be fixed once the go version of int works good enough for testing

the go version will be the only version going forward

im planning to release go at version 2.0

1.x is only python
2.x will only be go

version: 1.5.3

https://github.com/user-attachments/assets/219743bf-ece8-4439-bd79-cb15a3b970cc

# installation
## debian or debian based
download the .deb file
## other linux distros
copy int.py to any directory in your path

i recommend ~/.local/bin because it's usually not touched by package managers

## building the debian package
install git to download the repo ```sudo apt update && sudo apt install git```

clone the repo ```git clone https://github.com/WerIstLuka/int```

change directory to int ```cd int```

use the mkdpkg script to build ```bash mkdpkg```

