package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Address structure
type Address struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Email       string `json:"email,omitempty"`
}

var addresses map[string]Address
var addressFileName string

///////////////////////// DB Functions

// AddrID - create the key from last name and first name values
func AddrID(firstName string, lastName string) string {
	return strings.ToLower(lastName + "," + firstName)
}

// InitCache - init internal cache from file.
func InitCache(fromFileName string, fromJSON []byte) {
	addresses = make(map[string]Address)
	addressFileName = fromFileName
	if fromFileName != "" {
		fromJSON, _ = ioutil.ReadFile(addressFileName)
	}

	if string(fromJSON) == "" {
		// Using no backing store.
		return
	}

	var addressList []Address
	json.Unmarshal([]byte(fromJSON), &addressList)
	for _, address := range addressList {
		if _, ok := GetEntry(address.FirstName, address.LastName); !ok {
			addrID := AddrID(address.FirstName, address.LastName)
			addresses[addrID] = address
		} else {
			fmt.Printf("Initialization file contained a duplicate. %v\n", address)
		}
	}
}

// WriteToFile - update file
func WriteToFile() {
	if addressFileName == "" {
		// no backing store
		return
	}
	var addressList []Address
	for _, v := range addresses {
		addressList = append(addressList, v)
	}
	buffer, _ := json.MarshalIndent(addressList, "", "  ")
	err := ioutil.WriteFile(addressFileName, buffer, 0644)
	if err != nil {
		fmt.Printf("Error writing to backing store. %s\n", err.Error())
	}
}

// UpdateEntry - Update or add entry.
func UpdateEntry(address Address) {
	addresses[AddrID(address.FirstName, address.LastName)] = address
	WriteToFile()
}

// DeleteEntry - Delete a entry from cache
func DeleteEntry(address Address) {
	addrID := AddrID(address.FirstName, address.LastName)
	if _, ok := addresses[addrID]; ok {
		delete(addresses, addrID)
	}
	WriteToFile()
}

// GetEntry - Retrieve a entry from cache
func GetEntry(firstName string, lastName string) (Address, bool) {
	addrID := AddrID(firstName, lastName)
	if address, ok := addresses[addrID]; ok {
		return address, true
	}
	return Address{}, false
}

////////
/////////////////////////////////// Handlers
///////

// GetAddresses - Display all addresses
func GetAddresses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// get optional query values
	queryValues := r.URL.Query()
	formatSpec := queryValues.Get("format")
	if formatSpec != "" {
		if formatSpec != "csv" && formatSpec != "json" {
			http.Error(w, fmt.Sprintf("Invalid format specifier: %s", formatSpec), 500)
			return
		}
	}

	// lets sort the entries..
	if formatSpec == "csv" {
		csvBuffer := &bytes.Buffer{}
		csvWriter := csv.NewWriter(csvBuffer)

		for _, address := range addresses {
			var record []string

			record = append(record, address.FirstName)
			record = append(record, address.LastName)
			record = append(record, address.PhoneNumber)
			record = append(record, address.Email)
			csvWriter.Write(record)
		}
		csvWriter.Flush()
		w.Write(csvBuffer.Bytes())
	} else {
		var addressList []Address
		for _, address := range addresses {
			addressList = append(addressList, address)
		}
		json.NewEncoder(w).Encode(addressList)
	}
}

// GetAddrEntry - display a single address entry using query parameters  firstName & lastName.
func GetAddrEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// get query parameters
	queryValues := r.URL.Query()
	lastName := queryValues.Get("lastName")
	firstName := queryValues.Get("firstName")
	address, ok := GetEntry(firstName, lastName)
	if !ok {
		http.Error(w, fmt.Sprintf("Address entry not found for. firstName: %s, lastName %s", lastName, firstName), 404)
		return
	}
	json.NewEncoder(w).Encode(address)
}

// CreateAddrEntry - create a new address entry
func CreateAddrEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, ok := GetEntry(address.FirstName, address.LastName)
	if ok {
		http.Error(w, fmt.Sprintf("Duplicate entry for firstName: %s, lastName: %s", address.FirstName, address.LastName), 400)
		return
	}
	UpdateEntry(address)
	json.NewEncoder(w).Encode(address)
}

// DeleteAddrEntry - Delete a address entry
func DeleteAddrEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// get query parameters
	queryValues := r.URL.Query()
	lastName := queryValues.Get("lastName")
	firstName := queryValues.Get("firstName")
	address, ok := GetEntry(firstName, lastName)
	if !ok {
		http.Error(w, fmt.Sprintf("Entry not found for firstName: %s, lastName: %s", firstName, lastName), 404)
		return
	}
	DeleteEntry(address)

}

// UpdateAddrEntry - Update a address entry
func UpdateAddrEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, ok := GetEntry(address.FirstName, address.LastName)
	if !ok {
		http.Error(w, fmt.Sprintf("Entry not found. firstName: %s, lastName: %s", address.FirstName, address.LastName), 404)
		return
	}
	UpdateEntry(address)
}

// BulkImport - Will initialize and bulk import entries into the address book.
func BulkImport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// initialize the address book
	addresses = make(map[string]Address)

	var errors string

	csvReader := csv.NewReader(r.Body)
	csvReader.TrimLeadingSpace = true
	csvTable, err := csvReader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for _, record := range csvTable {
		var address Address
		address.FirstName = record[0]
		address.LastName = record[1]
		address.PhoneNumber = record[2]
		address.Email = record[3]
		if _, ok := GetEntry(address.FirstName, address.LastName); ok {
			errors += fmt.Sprintf("%v\n", address)
		} else {
			UpdateEntry(address)
		}
	}
	if errors != "" {
		http.Error(w, fmt.Sprintf("These entries were not processed as they were duplicate entries: %s\n", errors), 202)
	}
}

func InitRouterContext() *httprouter.Router {
	router := httprouter.New()

	// get list of people
	router.GET("/addresses", GetAddresses)

	// Bulk import addresses
	router.PUT("/addresses", BulkImport)

	// get specific address given first & last name
	router.GET("/addresses/address", GetAddrEntry)

	// add person
	router.POST("/addresses/address", CreateAddrEntry)

	// update address
	router.PUT("/addresses/address", UpdateAddrEntry)

	// delete address given first & last name
	router.DELETE("/addresses/address", DeleteAddrEntry)

	return router
}

func main() {
	fileName := flag.String("f", "", "[optional] Address file.")
	flag.Parse()

	InitCache(*fileName, []byte{})
	router := InitRouterContext()

	log.Fatal(http.ListenAndServe(":8000", router))

}
