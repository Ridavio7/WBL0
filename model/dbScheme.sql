create table delivery (
	id	bigserial not null primary key,
	Name    varchar(256),
	Phone   varchar(256),
	Zip     varchar(256), 
	City    varchar(256),
	Address varchar(256),
	Region  varchar(256),  
	Email   varchar(256)
);

create table payment (
	id	bigserial not null primary key,
	Transaction  varchar(256),
	RequestId    varchar(128),
	Currency     varchar(128), 
	Provider     varchar(128),
	Amount       int,
	PaymentDt    int,  
	Bank         varchar(128),
	DeliveryCost int,
	GoodsTotal   int,
	CustomFee    int
);

create table items (
	id	bigserial not null primary key,
	ChrtID     int,
	TrackNumber varchar(256),
	Price      int,    
	Rid        varchar(256), 
	Name       varchar(128), 
	Sale       int,    
	Size       varchar(128), 
	TotalPrice int,    
	NmID       int,    
	Brand      varchar(128),
	Status     int
);

create table "order" (
	id	bigserial not null primary key,
	OrderUID          varchar(128),
	TrackNumber       varchar(256),
	Entry             varchar(128),
	Locale            varchar(128), 
	InternalSignature varchar(128),
	CustomerID        varchar(128), 
	DeliveryService   varchar(128), 
	Shardkey          varchar(128),  
	SmID              int,
	DataCreated       varchar(256),
	OofShard          varchar(128)
);

create table cache (
	id	bigserial not null primary key,
	OrderUID          varchar(128),
	TrackNumber       varchar(256),
	Name    		  varchar(256),
	Phone   		  varchar(256),
	Zip     		  varchar(256), 
	City    		  varchar(256),
	Address 		  varchar(256),
	Region  		  varchar(256),  
	Email   		  varchar(256),
	ChrtID     		  int,
	TrackNumber_s 	  varchar(256),
	Price      		  int,    
	Rid        		  varchar(256), 
	Name_s       	  varchar(128), 
	Sale       		  int,    
	Size       		  varchar(128), 
	TotalPrice 		  int,    
	NmID       		  int,    
	Brand      		  varchar(128),
	Status     		  int
);