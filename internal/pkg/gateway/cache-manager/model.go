package cachemanager

type PartnersResponse struct {
	Data  []*Partner `json:"data,omitempty"`
	Error *Error     `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Address struct {
	CountryName  string `json:"countryName"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	AddressLine4 string `json:"addressLine4"`
	City         string `json:"city"`
	District     string `json:"district"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	CountryCode  string `json:"countryCode"`
	LanguageCode string `json:"languageCode"`
}

type Partner struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	DistributorID    string  `json:"distributorId"`
	Market           string  `json:"market"`
	Active           bool    `json:"active"`
	ParmaPartnerCode string  `json:"parmaPartnerCode"`
	RoleCode         string  `json:"roleCode"`
	Address          Address `json:"address"`
}
