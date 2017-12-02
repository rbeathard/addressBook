package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var testData string = `[
  {
    "firstName": "Daffy",
    "lastName": "Duck",
    "phoneNumber": "9085461001",
    "email": "daffy@google.com"
  },
  {
    "firstName": "Yosemite",
    "lastName": "Sam",
    "phoneNumber": "2143336666",
    "email": "yosemite@azure.com"
  },
  {
    "firstName": "Tweety",
    "lastName": "Bird",
    "phoneNumber": "7136448908",
    "email": "tweety@aws.com"
  },
  {
    "firstName": "Bugs",
    "lastName": "Bunny",
    "phoneNumber": "9725551212",
    "email": "bunny@google.com"
  }
]`

// getAddressesTest - very simple get all addresses..
func getAddressesTest(router *httprouter.Router, t *testing.T) {
	req, err := http.NewRequest("GET", "/addresses", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/addresses?format=csv", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
func getAddressTest(router *httprouter.Router, t *testing.T) {
	req, err := http.NewRequest("GET", "/addresses/address?firstName=Bugs&lastName=Bunny", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expected string
	expected = `{"firstName":"Bugs","lastName":"Bunny","phoneNumber":"9725551212","email":"bunny@google.com"}`
	expected += "\n"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func postAddressTest(router *httprouter.Router, t *testing.T) {
	entry := `{"firstName":"Foghorn","lastName":"LegHorn","phoneNumber":"9725551414","email":"foghorn@cisco.com"}`
	req, err := http.NewRequest("POST", "/addresses/address", strings.NewReader(entry))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/addresses/address?firstName=Foghorn&lastName=LegHorn", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	var expected string
	expected = entry + "\n"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func putAddressTest(router *httprouter.Router, t *testing.T) {
	entry := `{"firstName":"Foghorn","lastName":"LegHorn","phoneNumber":"408551212","email":"foghorn@cisco.com"}`
	req, err := http.NewRequest("PUT", "/addresses/address", strings.NewReader(entry))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/addresses/address?firstName=Foghorn&lastName=LegHorn", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	var expected string
	expected = entry + "\n"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func deleteAddressTest(router *httprouter.Router, t *testing.T) {
	req, err := http.NewRequest("DELETE", "/addresses/address?firstName=Foghorn&lastName=LegHorn", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/addresses/address?firstName=Foghorn&lastName=LegHorn", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func putAddressesTest(router *httprouter.Router, t *testing.T) {
	entry := `Geddy,Lee, 2143331213,r@hex.net
            Alex, Liefson,2143331213,r@hex.net
            Neil,Peart,2143331213,r@hex.net`
	req, err := http.NewRequest("PUT", "/addresses", strings.NewReader(entry))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHealthCheckHandler(t *testing.T) {

	InitCache("", []byte(testData))
	router := InitRouterContext()

	getAddressesTest(router, t)
	getAddressTest(router, t)
	postAddressTest(router, t)
	putAddressTest(router, t)
	deleteAddressTest(router, t)
	putAddressesTest(router, t)

}
