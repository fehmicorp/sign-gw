package config

import (
	"time"
)

type Ldap struct {
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	BaseDN   string `yaml:"baseDN"`
	BindDN   string `yaml:"bindDN"`
	Password string `yaml:"password"`
	UseTLS   bool   `yaml:"useTLS"`
}

type User struct {
	DN                string    `json:"dn"`
	CN                string    `json:"cn"`
	UserName          string    `json:"username"`
	UserPrincipalName string    `json:"userPrincipalName"`
	DisplayName       string    `json:"displayName"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Initials          string    `json:"initials"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Mobile            string    `json:"mobile"`
	Fax               string    `json:"fax,omitempty"`
	Office            string    `json:"office"`
	Description       string    `json:"description"`
	EmployeeID        string    `json:"employeeId"`
	EmployeeNumber    string    `json:"employeeNumber"`
	Department        string    `json:"department"`
	Division          string    `json:"division"`
	Title             string    `json:"title"`
	Company           string    `json:"company"`
	Manager           string    `json:"manager"`
	Street            string    `json:"street"`
	City              string    `json:"city"`
	State             string    `json:"state"`
	Country           string    `json:"country"`
	PostalCode        string    `json:"postalCode"`
	AccountEnabled    bool      `json:"accountEnabled"`
	AccountLocked     bool      `json:"accountLocked"`
	AccountExpired    bool      `json:"accountExpired"`
	LastLogon         time.Time `json:"lastLogon,omitempty"`
	PasswordLastSet   time.Time `json:"passwordLastSet,omitempty"`
	ObjectGUID        string    `json:"objectGuid"`
	ObjectSID         string    `json:"objectSid"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
	Groups            []string  `json:"groups"`
	ThumbnailPhoto    []byte    `json:"thumbnailPhoto,omitempty"`
}
