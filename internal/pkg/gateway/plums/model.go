package plums

type User struct {
	UserID         string         `json:"userId"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	Email          string         `json:"email"`
	CDSID          string         `json:"cdsid"`
	CountryCode    string         `json:"countryCode"`
	Partners       []Partner      `json:"partners"`
	UserIdentities []UserIdentity `json:"userIdentities"`
}

type Partner struct {
	PartnerID   string   `json:"partnerId"`
	PartnerType string   `json:"partnerType"` // PARMA, NSC, etc.
	IsPrimary   bool     `json:"isPrimary"`
	Roles       []string `json:"roles"`
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Users struct {
	HasNext bool    `json:"hasNext"`
	Result  []*User `json:"result"`
}

type UserIdentity struct {
	Provider       string `json:"provider"`
	ProviderUserID string `json:"providerUserId"`
	AccountName    string `json:"accountName"`
}
