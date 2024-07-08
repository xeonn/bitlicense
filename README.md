# Simple Software Licensing System

| :exclamation:  This project is in its very early stage   |
|----------------------------------------------------------|

## Overview
A Simple Go package for lightweight Software Licensing System designed to manage software licenses efficiently. It allows software vendors to control access to their software products through license keys, ensuring that only authorized users can utilize the software.

It uses PKI certificates generated from `ssh-keygen -t id25519` for hashing and signing license key.

## Feature
* License Key Generation: Generate unique license keys for each system signed by private key.
* License Validation: Validate license keys to ensure they are legitimate by validating the signature of the certificate.
* License Expiry: Set expiration dates for license keys.

### Planned
* License Management: Manage users/organization and their associated licenses.
* API Access: RESTful API for integrating the licensing system with other software.

## Requirement
* Golang v1.22.4 or higher

## Installation
Clone the repository into a folder named `bitlicense.git`
```
git clone https://github.com/xeonn/bitlicense.git bitlicense.git
```


## Usage

### License key generation
To generate a license key, run the following command
```
cd bitlicense.git/server
go run *.go generate --client demo --expiry 2023-06-05 > demo.json
```

a json file `demo.json` is created from output of the command above. 

### License key validation
To validate a license key, run the following command
```
cd bitlicense.git/server
go run *.go validate --file demo.json
```

### In your Go project
Import bitlicense into your project
```
go get github.com/xeonn/bitlicense
```
create a folder named `certs` and place your generated public key named `id_ed25519.pub` (this is hardcoded for now)
```
mkdir certs
cp [your_cert_file] certs/id_ed25519.pub
```
here's an example for using the validation
```
package main

import "github.com/xeonn/bitlicense"

func main() {
	if bitlicense.ValidateFile("demo.json") {
		println("License is VALID")
	} else {
		println("License is INVALID")
	}
}
``