<h1 align="center">
  <br>
  <a href="https://globaldce.github.io">
    <img src="globaldce.png" alt="globaldce" width="70%">
  </a>
  <br>
  globaldce gateway
  <br>
  <br>
</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/globaldce/globaldce-gateway)](https://goreportcard.com/report/github.com/globaldce/globaldce-gateway)

## globaldce-gateway
This is the reference implementation of globaldce protocol coded in the go programming language.

## What is globaldce ?
It is a decentralized applications platform built on a privacy focused peer-to-peer content delivery network powered by its own proof-of-work cryptocurrency. It features anonymous communication, censorship resistance, native registration, an integrated reputation system and version control capabilities.

## What is a globaldce gateway ?
A globaldce gateway acts as a bridge between users and globaldce. Through the gateway, users can handle their wallet assets while accessing and interacting with globaldce based decentralized applications as if they were traditional web server based applications.

## GUI Screenshots

<p float="left">
<img src="https://github.com/globaldce/globaldce-gateway/blob/main/screenshot1.jpg" width="45%">
<img src="https://github.com/globaldce/globaldce-gateway/blob/main/screenshot2.jpg" width="45%">
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
This software is open source, feel free to contribute! Not just code, but reporting issues, suggesting ideas, improving documentation or marketing material.

## Licence
This project is an open source free software; you can redistribute it and/or modify it under the terms of the MIT license.
See [LICENSE](https://github.com/globaldce/globaldce-gateway/blob/main/LICENSE) for details. 
