package nats

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgresql"
	dbname   = "postgres"
)

var clientIDSub = "SubID"

// подписка, получение данных от pub, отправка в БД
func Sub() {
	var (
		URL       string
		userCreds string
		qgroup    string
		durable   string
	)

	opts := []nats.Option{nats.Name("NATS Streaming Example Subscriber")}
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientIDSub, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to clusterID: [%s] clientID: [%s]\n", clusterID, clientIDSub)

	startOpt := stan.StartWithLastReceived() // берет последнее сообзение
	i := 0
	mcb := func(msg *stan.Msg) {
		i++
		printMsg(msg, i)       // постит в консоль полученные данне
		validAndSend(msg.Data) // валидация и отправка в БД
	}

	_, err = sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clusterID=[%s], clientID=[%s]\n", subj, clusterID, clientIDSub)
}

func printMsg(m *stan.Msg, i int) {
	log.Printf("Get message [#%d]: %s\n", i, m)
}

func validAndSend(data []byte) {
	var r Order
	validJSON := json.Valid(data) //проверяем полученные данные
	if validJSON == true {
		err := json.Unmarshal(data, &r) //декодируем данные
		if err != nil {
			log.Printf("Error, %v\n", err)
		}
		log.Printf("Unmarshal data: %v\n", r)
		log.Printf("Sending data...")
		sendToDB(r)                              // отправляем в БД
		Ccache.Set(r.OrderUID, r, 5*time.Minute) // создаем cache
		log.Printf("Caching data...")
		log.Printf("Cache is living fo 5 minutes!")
	}
}

func sendToDB(r Order) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	defer db.Close()

	// заполняем таблицу delivery
	_, err = db.Exec(`insert into "delivery" ("name", "phone", "zip", "city", "address", "region", "email") values ($1, $2, $3, $4, $5, $6, $7)`,
		r.Delivery.Name, r.Delivery.Phone, r.Delivery.Zip, r.Delivery.City, r.Delivery.Address, r.Delivery.Region, r.Delivery.Email)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	// заполняем таблицу payment
	_, err = db.Exec(`insert into "payment" ("transaction", "requestid", "currency", "provider", "amount", "paymentdt", "bank", "deliverycost", "goodstotal", "customfee") 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, r.Payment.Transaction, r.Payment.RequestId, r.Payment.Currency, r.Payment.Provider, r.Payment.Amount, r.Payment.PaymentDt, r.Payment.Bank,
		r.Payment.DeliveryCost, r.Payment.GoodsTotal, r.Payment.CustomFee)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	// заполняем таблицу items
	_, err = db.Exec(`insert into "items" ("chrtid", "tracknumber", "price", "rid", "name", "sale", "size", "totalprice", "nmid", "brand", "status") 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, r.Items.ChrtID, r.Items.TrackNumber, r.Items.Price, r.Items.Rid, r.Items.Name, r.Items.Sale,
		r.Items.Size, r.Items.TotalPrice, r.Items.NmID, r.Items.Brand, r.Items.Status)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	// заполняем таблицу order
	_, err = db.Exec(`insert into "order" ("orderuid", "tracknumber", "entry", "locale", "internalsignature", "customerid", "deliveryservice", "shardkey", "smid", "datacreated", "oofshard") 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, r.OrderUID, r.TrackNumber, r.Entry, r.Locale, r.InternalSignature, r.CustomerID, r.DeliveryService,
		r.Shardkey, r.SmID, r.DataCreated, r.OofShard)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	// заполняем таблицу для cache
	_, err = db.Exec(`insert into "cache" ("orderuid", "tracknumber", "name", "phone", "zip", "city", "address", "region", "email",
	"chrtid", "tracknumber_s", "price", "rid", "name_s", "sale", "size", "totalprice", "nmid", "brand", "status") 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`, r.OrderUID, r.TrackNumber, r.Delivery.Name, r.Delivery.Phone,
		r.Delivery.Zip, r.Delivery.City, r.Delivery.Address, r.Delivery.Region, r.Delivery.Email, r.Items.ChrtID, r.Items.TrackNumber,
		r.Items.Price, r.Items.Rid, r.Items.Name, r.Items.Sale, r.Items.Size, r.Items.TotalPrice, r.Items.NmID, r.Items.Brand, r.Items.Status)
	if err != nil {
		log.Printf("Error %s\n", err)
	}

	log.Printf("Data is allready in DB!")

}
