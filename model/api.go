package model

import (
	"database/sql"
	"fmt"
	"log"

	"html/template"
	"net/http"

	n "WBL0/nats"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/mitchellh/mapstructure"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgresql"
	dbname   = "postgres"
)

var CacheValue n.Order
var ShowOrder n.Order

func InitServer() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", routIndex).Methods("GET")
	rtr.HandleFunc("/returnId", returnId).Methods("POST")

	http.Handle("/", rtr)

	log.Println("Start server http://localhost:5555")
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Printf("Error %s\n", err)
	}
}

func routIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")

	if err != nil {
		log.Printf("Error %s\n", err)
	}

	t.ExecuteTemplate(w, "index", ShowOrder)
}

func returnId(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	defer db.Close()

	rows, err := db.Query(`SELECT "orderuid", "tracknumber", "name", "phone", "zip", "city", "address", "region", "email",
	"chrtid", "tracknumber", "price", "rid", "name", "sale", "size", "totalprice", "nmid", "brand", "status" FROM cache`)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c n.Order
		err = rows.Scan(&c.OrderUID, &c.TrackNumber, &c.Delivery.Name, &c.Delivery.Phone, &c.Delivery.Zip, &c.Delivery.City, &c.Delivery.Address,
			&c.Delivery.Region, &c.Delivery.Email, &c.Items.ChrtID, &c.Items.TrackNumber, &c.Items.Price, &c.Items.Rid, &c.Items.Name, &c.Items.Sale,
			&c.Items.Size, &c.Items.TotalPrice, &c.Items.NmID, &c.Items.Brand, &c.Items.Status)
		if err != nil {
			log.Printf("Error %s\n", err)
		}

		CacheValue = c

	}

	if key == CacheValue.OrderUID {
		value, _ := n.Ccache.Get(key)
		log.Printf("Key is valid!")
		c := n.Order{}
		mapstructure.Decode(value, &c)

		ShowOrder = c

	} else {
		log.Printf("Key isn`t valid!")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
