## globaldce-toolbox
This is the reference implementation of globaldce protocole coded in the go programming language.

globaldce is a blockchain based experimental decentralized collaboration environnement. While this project already delivers some of the features we want to see, it is still at an early stage and many features needs to be implemented and tested.

## Getting started
We recommend that new users get started using the compressed executable binaries that can be found in releases section:
https://github.com/globaldce/globaldce-toolbox/releases

Just download the compressed release prepared for your operating system, decompress it and run the globaldce-toolbox executable file.

## Current Status
+ Tested on Go 1.5
+ Tested on Windows, Linux and MacOSX

## How to compile from source
1. If you have not already done so, install Go (at least version 1.12)
2. As this code includes a graphical user interface, you will need to install a C compiler (to handle graphics drivers) and graphics drivers. More information on how to do this, depending on your operating system, can be found here:
https://developer.fyne.io/started/
3. Use the following commands to download, and build this source code: 
```bash
git clone https://github.com/globaldce/globaldce-toolbox
cd globaldce-toolbox
go build -o globaldce-toolbox
```

## Usage
This toolbox cli already offers some of the project features. To get the available commands run the node as follow:
```bash
globaldce-toolbox -help
```

## Contributing
Bug reports and bug fixes are welcome.

## Licence
This project is an open source free software; you can redistribute it and/or modify it under the terms of the MIT license.
See [LICENSE](https://github.com/globaldce/globaldce-toolbox/blob/main/LICENSE) for details. 
