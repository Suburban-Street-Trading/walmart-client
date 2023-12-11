package walmart

import "github.com/lomsa-dev/gonull"

type OrderStatus string

const (
	OrderStatusCreated      OrderStatus = "Created"
	OrderStatusAcknowledged OrderStatus = "Acknowledged"
	OrderStatusShipped      OrderStatus = "Shipped"
	OrderStatusDelivered    OrderStatus = "Delivered"
	OrderStatusCancelled    OrderStatus = "Cancelled"
	OrderStatusRefund       OrderStatus = "Refund"
)

type MethodCode string

const (
	MethodCodeStandard MethodCode = "Standard"
	MethodCodeExpress  MethodCode = "Express"
	MethodCodeValue    MethodCode = "Value"
)

type Credentials struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
}

// AUTH TYPES

type AccessTokenResponse struct {
	AccessToken  string                  `json:"access_token"`
	TokenType    string                  `json:"token_type"`
	ExpiresIn    int                     `json:"expires_in"`
	RefreshToken gonull.Nullable[string] `json:"refresh_token"`
}

// GENERAL TYPES - these are used in multiple places

type Quantity struct {
	UnitOfMeasurement string `json:"unitOfMeasurement"`
	Amount            string `json:"amount"`
}

type Charge struct {
	ChargeType   string             `json:"chargeType"`
	ChargeName   string             `json:"chargeName"`
	ChargeAmount Amount             `json:"chargeAmount"`
	Tax          OrderLineChargeTax `json:"tax"`
}

type Amount struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// SHIPMENT TYPES

type ShipOrderLinesRequest struct {
	OrderShipment OrderShipment `json:"orderShipment"`
}

type OrderShipment struct {
	ProcessMode gonull.Nullable[string] `json:"processMode"`
	OrderLines  OrderShipmentLines      `json:"orderLines"`
}

type OrderShipmentLines struct {
	OrderLine []ShipmentOrderLine `json:"orderLine"`
}

type ShipmentOrderLine struct {
	LineNumber             string                `json:"lineNumber"`
	IntentToCancelOverride gonull.Nullable[bool] `json:"intentToCancelOverride"`
	SellerOrderID          string                `json:"sellerOrderId"`
	OrderLineStatuses      OrderLineStatuses     `json:"orderLineStatuses"`
}

type SupportedCarriersResponse struct {
	Carriers []Carrier `json:"carriers"`
	Errors   []any     `json:"errors"`
}

type Carrier struct {
	CarrierID        string `json:"carrierId"`
	CarrierShortName string `json:"carrierShortName"`
	CarrierName      string `json:"carrierName"`
}

// ORDER TYPES

type GetOrdersResponse struct {
	List struct {
		Errors []any `json:"errors"`
		//Meta     Meta     `json:"meta"`
		Elements Elements `json:"elements"`
	} `json:"list"`
}

type Meta struct {
	TotalCount int    `json:"totalCount"`
	Limit      int    `json:"limit"`
	NextCursor string `json:"nextCursor"`
}

type Elements struct {
	Order []Order `json:"order"`
}

type Order struct {
	PurchaseOrderID         string                        `json:"purchaseOrderId"`
	CustomerOrderID         string                        `json:"customerOrderId"`
	CustomerEmailID         string                        `json:"customerEmailId"`
	OrderType               gonull.Nullable[string]       `json:"orderType"`
	OriginalCustomerOrderID gonull.Nullable[string]       `json:"originalCustomerOrderId"`
	OrderDate               int64                         `json:"orderDate"`
	BuyerID                 gonull.Nullable[string]       `json:"buyerId"`
	Mart                    gonull.Nullable[string]       `json:"mart"`
	IsGuest                 gonull.Nullable[bool]         `json:"isGuest"`
	ShippingInfo            OrderShippingInfo             `json:"shippingInfo"`
	OrderLines              OrderLines                    `json:"orderLines"`
	PaymentTypes            gonull.Nullable[[]string]     `json:"paymentTypes"`
	OrderSummary            gonull.Nullable[OrderSummary] `json:"orderSummary"`
	PickupPersons           any                           `json:"pickupPersons"`
	ShipNode                any                           `json:"shipNode"`
}

type OrderSummary struct {
	TotalAmount    TotalAmount `json:"totalAmount"`
	OrderSubTotals any         `json:"orderSubTotals"`
}

type TotalAmount struct {
	CurrencyAmount float64 `json:"currencyAmount"`
	CurrencyUnit   string  `json:"currencyUnit"`
}

type OrderShippingInfo struct {
	Phone                 string                   `json:"phone"`
	EstimatedDeliveryDate any                      `json:"estimatedDeliveryDate"`
	EstimatedShipDate     any                      `json:"estimatedShipDate"`
	MethodCode            string                   `json:"methodCode"`
	PostalAddress         OrderShippingInfoAddress `json:"postalAddress"`
}

type OrderLines struct {
	OrderLine []OrderLine `json:"orderLine"`
}

type OrderLine struct {
	LineNumber            string                       `json:"lineNumber"`
	Item                  OrderLineItem                `json:"item"`
	Charges               OrderLineCharges             `json:"charges"`
	OrderLineQuantity     Quantity                     `json:"orderLineQuantity"`
	StatusDate            int64                        `json:"statusDate"`
	OrderLineStatuses     OrderLineStatuses            `json:"orderLineStatuses"`
	ReturnOrderID         gonull.Nullable[string]      `json:"returnOrderId"`
	Refund                gonull.Nullable[Refund]      `json:"refund"`
	OriginalCarrierMethod gonull.Nullable[string]      `json:"originalCarrierMethod"`
	ReferenceLineID       gonull.Nullable[string]      `json:"referenceLineId"`
	Fulfillment           gonull.Nullable[Fulfillment] `json:"fulfillment"`
	SerialNumbers         gonull.Nullable[[]string]    `json:"serialNumbers"`
	IntentToCancel        gonull.Nullable[string]      `json:"intentToCancel"`
	ConfigID              gonull.Nullable[string]      `json:"configId"`
	SellerOrderID         gonull.Nullable[string]      `json:"sellerOrderId"`
}

type Fulfillment struct {
	FulfillmentOption   gonull.Nullable[string] `json:"fulfillmentOption"`
	ShipMethod          gonull.Nullable[string] `json:"shipMethod"`
	StoreID             gonull.Nullable[string] `json:"storeId"`
	PickUpDateTime      gonull.Nullable[int64]  `json:"pickUpDateTime"`
	PuckUpBy            gonull.Nullable[string] `json:"pickUpBy"`
	ShippingProgramType gonull.Nullable[string] `json:"shippingProgramType"`
}

type Refund struct {
	RefundID      gonull.Nullable[string] `json:"refundId"`
	RefundComment gonull.Nullable[string] `json:"refundComment"`
	RefundCharges RefundCharges           `json:"refundCharges"`
}

type RefundCharges struct {
	RefundCharge []RefundCharge `json:"refundCharge"`
}

type RefundCharge struct {
	RefundReason string `json:"refundReason"`
	Charge       Amount `json:"charge"`
}

type OrderLineStatuses struct {
	OrderLineStatus []OrderLineStatus `json:"orderLineStatus"`
}

type OrderLineStatus struct {
	Status             OrderStatus             `json:"status"`
	StatusQuantity     Quantity                `json:"statusQuantity"`
	CancellationReason gonull.Nullable[string] `json:"cancellationReason"`
	TrackingInfo       TrackingInfo            `json:"trackingInfo"`
}

type TrackingInfo struct {
	ShipDateTime   int64                   `json:"shipDateTime"`
	CarrierName    CarrierName             `json:"carrierName"`
	MethodCode     MethodCode              `json:"methodCode"`
	TrackingNumber string                  `json:"trackingNumber"`
	TrackingURL    gonull.Nullable[string] `json:"trackingURL"`
}

type CarrierName struct {
	OtherCarrier gonull.Nullable[string] `json:"otherCarrier"`
	Carrier      gonull.Nullable[string] `json:"carrier"`
}

type OrderLineCharges struct {
	OrderLineCharge []Charge `json:"charge"`
}

type OrderLineChargeTax struct {
	TaxName   string `json:"taxName"`
	TaxAmount Amount `json:"taxAmount"`
}

type OrderLineItem struct {
	ProductName string                  `json:"productName"`
	SKU         string                  `json:"sku"`
	ImageURL    gonull.Nullable[string] `json:"imageUrl"`
	Weight      OrderLineItemWeight     `json:"weight"`
}

type OrderLineItemWeight struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type OrderShippingInfoAddress struct {
	Name        string                  `json:"name"`
	Address1    string                  `json:"address1"`
	Address2    gonull.Nullable[string] `json:"address2"`
	City        string                  `json:"city"`
	State       string                  `json:"state"`
	PostalCode  string                  `json:"postalCode"`
	Country     string                  `json:"country"`
	AddressType gonull.Nullable[string] `json:"addressType"`
}
