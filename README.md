# Address Book Sample Code in Golang
The following repo contains example golang code based on the following requirements:

#### Requirements

As a user, I need an online address book exposed as a REST API.  I need the data set to include the following data fields:
First Name, Last Name, Email Address, and Phone Number

I need the api to follow standard rest semantics to support listing entries, showing a specific single entry, and adding, modifying, and deleting entries.  

The code for the address book should include regular go test files that demonstrate how to exercise all operations of the service.  

Finally I need the service to provide endpoints that can export and import the address book data in a CSV format.  

#### To Build
_Note: The following program utilizes **github.com/julienschmidt/httprouter** package and has been included in the package via the vendor directory._

Set GOPATH variable to working directory
> export GOPATH=$PWD

Get package
> go get -u github.com/rbeathard/addressBook

To run test
> go test -cover github.com/rbeathard/addressBook

#### Run program
addressBook can be ran in one of two modes: in memory or disk backed mode. In the in memory mode, any updates will not lost upon termination of program. In disk backed mode, updates will be written to disk. The optional -f flag specifies the disk file.  address.json has been included to test.

To start addressBook using the address.json test file. addressBook will listen on port 8000.
>  bin/addressBook -f src/github.com/rbeathard/addressBook/address.json

To list out the addresses in the address book, paste the following in a browser.
> http://localhost:8000/addresses



#### Endpoint supported

* *GET /addresses* - return all address book entries
* *PUT /addresses* - bulk upload addresses via csv format. Operation will overwrite all existing addresses
* *GET /addresses/address* - get a address given first and last name specified via query parameters firstName and lastName.
* *POST /addresses/address* - add a address entry
* *PUT /addresses/address* - update an existing entry
* *DELETE /addresses/address* - delete a address entry given first and last name specified via query parameters firstName and lastName.
