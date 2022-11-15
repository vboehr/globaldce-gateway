## globaldce-gateway
This is the reference implementation of globaldce protocole coded in the go programming language.

## What is globaldce ?
Even though cryptocurrency and blockchain technologies have the potential to disrupt every aspect of our daily life their use cases are currently still limited. Most real-world situations can not be tackled using existing smart contracts and dapps platforms and this due to the fact that building smart contracts and decentralized applications using existing platforms is like trying to building skyskrapers on sands.


Existing smart contracts and dapps platforms do not have built-in native decentralization, governance, privacy and scalability. And that led to painful and catastrophic situations for users.
* A lack of decentralization threatens the uncensored nature of blockchains
* A lack of governance leads to numerous forking
* A lack of privacy makes linking users addresses to their real-world identities easy
* A lack of scalability slows transactions and increases fees

## What is a globaldce gateway ?
A globaldce gateway acts as a bridge between users and globaldce. Through the gateway, users can handle their wallet assets while accessing and interacting with globaldce based decentralized applications as if they were traditional web server based applications.

## GUI Screenshots

<p float="left">
<img src="https://github.com/globaldce/globaldce-gateway/blob/main/screenshot1.png" width="45%">
<img src="https://github.com/globaldce/globaldce-gateway/blob/main/screenshot2.png" width="45%">
</p>

## Getting started
We recommend that new users get started using the compressed executable binaries that can be found in releases section:
https://github.com/globaldce/globaldce-gateway/releases

Just download the compressed release prepared for your operating system, decompress it and run the globaldce-gateway executable file.

## Current Status
+ Tested on Go 1.5
+ Tested on Windows, Linux and MacOSX

## How to compile from source
1. If you have not already done so, install Go (at least version 1.12)
2. As this code includes a graphical user interface, you will need to install a C compiler (to handle graphics drivers) and graphics drivers. More information on how to do this, depending on your operating system, can be found here:
https://developer.fyne.io/started/
3. Use the following commands to download, and build this source code: 
```bash
git clone https://github.com/globaldce/globaldce-gateway
cd globaldce-gateway
go build -o globaldce-gateway
```

## Usage
This code includes a graphical user interface (gui) and a command line interface (cli). 

In order to run the graphical user interface just run the node as follow:
```bash
globaldce-gateway
```

This cli already offers some of the project features. To get the available commands run as follow:
```bash
globaldce-gateway -help
```

## Contributing
Bug reports and bug fixes are welcome.

## Licence
This project is an open source free software; you can redistribute it and/or modify it under the terms of the MIT license.
See [LICENSE](https://github.com/globaldce/globaldce-gateway/blob/main/LICENSE) for details. 
