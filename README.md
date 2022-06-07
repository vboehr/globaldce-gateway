## globaldce-gateway
This is the reference implementation of globaldce protocole coded in the go programming language.

## What is globaldce ?
Even though cryptocurrency and blockchain technologies have the potential to disrubt every aspect of our daily lifes, their use cases are currently still limited to financial services. Most real-world situations can not be tackled using existing dapps platforms and this due to the fact they are missing a crucial aspect: nondeterministic functionalities.

Up until this point innovation in nondeterministic decentralized applications such as storage and datafeeds have been highly fragmented and implementing any of these technologies has required creating an entire protocol layer or even a specialized blockchain. However, these innovations can be made an order of magnitude easier to implement if only there was a broader foundational layer for all of these dapps to be built up. And this need is what we seek to satisfy.

While other blockchains only provide deterministic functionalities and force teams to rely on third party services for nondeteministic aspects, we believe here that blockchain based nondeterministic features coupled with smart contract functionalities can bring about a unprecedented level of collaboration and unleash the full potential of web 3.

## What is a globaldce gateway ?
A globaldce gateway acts as a bridge between users and globaldce. Through the gateway, users can handle their wallet assets while accessing and interacting with globaldce based decentralized applications as if they were traditional web server based applications.

## Screenshots

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
Bug reports and bug fixes are welcome.

## Licence
This project is an open source free software; you can redistribute it and/or modify it under the terms of the MIT license.
See [LICENSE](https://github.com/globaldce/globaldce-gateway/blob/main/LICENSE) for details. 
