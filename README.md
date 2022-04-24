## go-globaldce
This is the reference implementation of globaldce protocole coded in the go programming language.

## What is globaldce ?
globaldce is a blockchain based experimental decentralized collaboration environment. The goal of of the globaldce is to provide a fullstack collaborative environment for decentralized applications. While other web3.0 blockchains only provide backend smart contract functionalities and force teams to rely on web2.0 servers for their front end, we believe here that in order to unleash the full potential of web3.0, decentralized applications should not rely on centralized servers. Furthermore, we believe also that a blockchain based decentralized front end environment coupled with smart contract functionalities can bring about a unprecedented level of collaboration.

While this project already delivers some of the features we want to see, it is still at an early stage and many features needs to be implemented and tested.


## What is a globaldce gateway ?
A globaldce gateway acts as a bridge between users and globaldce. Through the gateway, users can handle their wallet assets while accessing and interacting with globaldce based decentralized applications as if they were traditional web server based applications.


## Getting started
We recommend that new users get started using the compressed executable binaries that can be found in releases section:
https://github.com/globaldce/go-globaldce/releases

Just download the compressed release prepared for your operating system, decompress it and run the go-globaldce executable file.

## Current Status
+ Tested on Go 1.5
+ Tested on Windows, Linux and MacOSX

## How to compile from source
1. If you have not already done so, install Go (at least version 1.12)
2. As this code includes a graphical user interface, you will need to install a C compiler (to handle graphics drivers) and graphics drivers. More information on how to do this, depending on your operating system, can be found here:
https://developer.fyne.io/started/
3. Use the following commands to download, and build this source code: 
```bash
git clone https://github.com/globaldce/go-globaldce
cd go-globaldce
go build -o go-globaldce
```

## Usage
This code includes a graphical user interface (gui) and a command line interface (cli). 

In order to run the graphical user interface just run the node as follow:
```bash
go-globaldce
```

This cli already offers some of the project features. To get the available commands run as follow:
```bash
go-globaldce -help
```

## Contributing
Bug reports and bug fixes are welcome.

## Licence
This project is an open source free software; you can redistribute it and/or modify it under the terms of the MIT license.
See [LICENSE](https://github.com/globaldce/go-globaldce/blob/main/LICENSE) for details. 
