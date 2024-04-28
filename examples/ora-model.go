package model

import (
	"time"
)

type CUSTOMERS struct {
	CustomerID interface{} `gorm:"primary_key;index:CUSTOMERS_PK;not null"`
	EmailAddress interface{} `gorm:"index:CUSTOMERS_EMAIL_U;not null"`
	FullName interface{} `gorm:"not null"`
}

func (c *CUSTOMERS) TableName() string { return "CUSTOMERS" }

type INVENTORY struct {
	InventoryID interface{} `gorm:"primary_key;index:INVENTORY_PK;not null"`
	StoreID interface{} `gorm:"index:INVENTORY_STORE_PRODUCT_U;not null"`
	ProductID interface{} `gorm:"index:INVENTORY_STORE_PRODUCT_U;index:INVENTORY_PRODUCT_ID_I;not null"`
	ProductInventory interface{} `gorm:"not null"`
	STORES STORES `gorm:"foreignKey:StoreID;references:StoreID"`
	PRODUCTS PRODUCTS `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (c *INVENTORY) TableName() string { return "INVENTORY" }

type ORDERS struct {
	OrderID interface{} `gorm:"primary_key;index:ORDERS_PK;not null"`
	OrderTms interface{} `gorm:"not null"`
	CustomerID interface{} `gorm:"index:ORDERS_CUSTOMER_ID_I;not null"`
	OrderStatus interface{} `gorm:"not null"`
	StoreID interface{} `gorm:"index:ORDERS_STORE_ID_I;not null"`
	CUSTOMERS CUSTOMERS `gorm:"foreignKey:CustomerID;references:CustomerID"`
	STORES STORES `gorm:"foreignKey:StoreID;references:StoreID"`
}

func (c *ORDERS) TableName() string { return "ORDERS" }

type ORDERITEMS struct {
	OrderID interface{} `gorm:"index:ORDER_ITEMS_PRODUCT_U;index:ORDER_ITEMS_PK;not null"`
	LineItemID interface{} `gorm:"index:ORDER_ITEMS_PK;not null"`
	ProductID interface{} `gorm:"index:ORDER_ITEMS_PRODUCT_U;not null"`
	UnitPrice interface{} `gorm:"not null"`
	Quantity interface{} `gorm:"not null"`
	ShipmentID interface{} `gorm:"index:ORDER_ITEMS_SHIPMENT_ID_I;not null"`
	ORDERS ORDERS `gorm:"foreignKey:OrderID;references:OrderID"`
	PRODUCTS PRODUCTS `gorm:"foreignKey:ProductID;references:ProductID"`
	SHIPMENTS SHIPMENTS `gorm:"foreignKey:ShipmentID;references:ShipmentID"`
}

func (c *ORDERITEMS) TableName() string { return "ORDER_ITEMS" }

type PRODUCTS struct {
	ProductID interface{} `gorm:"primary_key;index:PRODUCTS_PK;not null"`
	ProductName interface{} `gorm:"not null"`
	UnitPrice interface{} `gorm:"not null"`
	ProductDetails interface{} `gorm:"not null"`
	ProductImage interface{} `gorm:"not null"`
	ImageMimeType interface{} `gorm:"not null"`
	ImageFilename interface{} `gorm:"not null"`
	ImageCharset interface{} `gorm:"not null"`
	ImageLastUpdated time.Time `gorm:"not null"`
}

func (c *PRODUCTS) TableName() string { return "PRODUCTS" }

type SHIPMENTS struct {
	ShipmentID interface{} `gorm:"primary_key;index:SHIPMENTS_PK;not null"`
	StoreID interface{} `gorm:"index:SHIPMENTS_STORE_ID_I;not null"`
	CustomerID interface{} `gorm:"index:SHIPMENTS_CUSTOMER_ID_I;not null"`
	DeliveryAddress interface{} `gorm:"not null"`
	ShipmentStatus interface{} `gorm:"not null"`
	STORES STORES `gorm:"foreignKey:StoreID;references:StoreID"`
	CUSTOMERS CUSTOMERS `gorm:"foreignKey:CustomerID;references:CustomerID"`
}

func (c *SHIPMENTS) TableName() string { return "SHIPMENTS" }

type STORES struct {
	StoreID interface{} `gorm:"primary_key;index:STORES_PK;not null"`
	StoreName interface{} `gorm:"index:STORE_NAME_U;not null"`
	WebAddress interface{} `gorm:"not null"`
	PhysicalAddress interface{} `gorm:"not null"`
	Latitude interface{} `gorm:"not null"`
	Longitude interface{} `gorm:"not null"`
	Logo interface{} `gorm:"not null"`
	LogoMimeType interface{} `gorm:"not null"`
	LogoFilename interface{} `gorm:"not null"`
	LogoCharset interface{} `gorm:"not null"`
	LogoLastUpdated time.Time `gorm:"not null"`
}

func (c *STORES) TableName() string { return "STORES" }

