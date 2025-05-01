package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

func main() {
	var err error = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripe.Key = os.Getenv("SECRET_KEY")

	fmt.Println("Hello from go")

	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)
	http.HandleFunc("/health", handleHealth)

	log.Println("Listening on localhost:4242...")
	err = http.ListenAndServe("localhost:4242", nil)
	if err != nil {
		log.Fatal(err)
	}

	// var err error = returnsError("wrongpassword")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

// func returnsError(password string) error {
// 	var secretPassword string = "supersecretpassword"
// 	if password == secretPassword {
// 		return nil
// 	} else {
// 		return errors.New("invalid password")
// 	}
// }

// in go, must send response in slice of bytes (bytes array)
func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENDPOINT CALLED!")
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("request method was correct!")

	var req struct {
		ProductId string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address1  string `json:"address_1"`
		Address2  string `json:"address_2"`
		City      string `json:"city"`
		State     string `json:"state"`
		Zip       string `json:"zip"`
		Country   string `json:"country"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(req.ProductId)
	fmt.Println(req.FirstName)
	fmt.Println(req.LastName)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(req.ProductId)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
		// 	Enabled: stripe.Bool(true),
		// },
	}

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(paymentIntent.ClientSecret)

	// send client secret to stripe
	var response struct {
		ClientSecret string `json:"clientSecret"`
	}

	response.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//add Header to specify response is a json object
	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, &buf)
	if err != nil {
		fmt.Println(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// var response []byte = []byte("Server is up and running!")
	//byte type is a utf-8 string
	response := []byte("Server is up and running!")

	_, err := w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}

func calculateOrderAmount(productId string) int64 {
	switch productId {
	case "Forever Pants":
		return 26000
	case "Forever Shirt":
		return 15500
	case "Forever Shorts":
		return 30000
	}
	return 0
}

// in go, can pass funcs as args
// functionOne(anotherFunction)
// func functionOne(functionTwo func()) {
// 	functionTwo()
// }

// func anotherFunction() {
// 	fmt.Println("anotherFunction was called")
// }

//arrays in go
// var names []string = []string{"James", "Bill", "Jack"}
// var names []int = []int{23, 69, 420}

//returning multiple types
// func returnsMultiple() (string, int, bool){
// 	return "string", 1, true
// }

// shorthand. declare and assign var
// num := 1

// shorthand. go infers var type
// someString, someInt, someBool := returnsMultiple()
// fmt.Println(someString)
