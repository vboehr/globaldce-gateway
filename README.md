## globaldce-gateway
This is the reference implementation of globaldce protocole coded in the go programming language.

## What is globaldce ?
globaldce is a blockchain based experimental decentralized collaboration environment. The goal of of the globaldce is to provide a fullstack collaborative environment for decentralized applications. While other web3.0 blockchains only provide backend smart contract functionalities and force teams to rely on web2.0 servers for their front end, we believe here that in order to unleash the full potential of web3.0, decentralized applications should not rely on centralized servers. Furthermore, we believe also that a blockchain based decentralized front end environment coupled with smart contract functionalities can bring about a unprecedented level of collaboration.

## What is a globaldce gateway ?
A globaldce gateway acts as a bridge between users and globaldce. Through the gateway, users can handle their wallet assets while accessing and interacting with globaldce based decentralized applications as if they were traditional web server based applications.

## Screenshots

<p align="center">
  <img src="https://github.com/globaldce/globaldce-gateway/blob/main/screenshot1.jpg" width="50%" alt="screenshot">
  <img src="https://github.com/globaldce/globaldce-gateway/blob/main/screenshot2.jpg" width="50%" alt="screenshot">
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
