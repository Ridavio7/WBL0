package nats

import (
	"encoding/json"
	"log"

	_ "github.com/lib/pq"

	"github.com/nats-io/stan.go"
)

var clusterID = "test-cluster"
var clientIDPub = "PubID"
var URL = "nats://localhost:4222"
var subj = "order"

func Pub() {
	sc, err := stan.Connect(clusterID, clientIDPub, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at %s", err, URL)
	}
	defer sc.Close()

	delivery := Delivery{Name: "Test Testov", Phone: "+9720000000", Zip: "2639809", City: "Kiryat Mozkin", Address: "Ploshad Mira 15", Region: "Kraiot", Email: "test@gmail.com"}
	payment := Payment{Transaction: "b563feb7b2b84b6test", RequestId: "", Currency: "USD", Provider: "wbpay", Amount: 1817, PaymentDt: 1637907727, Bank: "alpha", DeliveryCost: 1500, GoodsTotal: 317, CustomFee: 0}
	items := Items{ChrtID: 9934930, TrackNumber: "WBILMTESTTRACK", Price: 453, Rid: "ab4219087a764ae0btest", Name: "Mascaras", Sale: 30, Size: "0", TotalPrice: 317, NmID: 2389212, Brand: "Vivienne Sabo", Status: 202}

	order := Order{
		OrderUID:          "b563feb7b2b84b6test",
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DataCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}
	msg, err := json.Marshal(order)

	if err != nil {
		log.Fatal(err)
	}

	sc.Publish(subj, msg)

	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
	}
	log.Printf("Published [%s] : '%s'\n", subj, msg)
}
