package ldap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	gldap "github.com/go-ldap/ldap/v3"
)

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

func entryToUser(e *gldap.Entry) config.User {

	uac, _ := strconv.Atoi(e.GetAttributeValue("userAccountControl"))

	return config.User{
		// Identity
		DN:                e.GetAttributeValue("distinguishedName"),
		CN:                e.GetAttributeValue("cn"),
		UserName:          e.GetAttributeValue("sAMAccountName"),
		UserPrincipalName: e.GetAttributeValue("userPrincipalName"),
		DisplayName:       e.GetAttributeValue("displayName"),
		FirstName:         e.GetAttributeValue("givenName"),
		LastName:          e.GetAttributeValue("sn"),
		Initials:          e.GetAttributeValue("initials"),

		// Contact
		Email:       e.GetAttributeValue("mail"),
		Phone:       e.GetAttributeValue("telephoneNumber"),
		Mobile:      e.GetAttributeValue("mobile"),
		Fax:         e.GetAttributeValue("facsimileTelephoneNumber"),
		Office:      e.GetAttributeValue("physicalDeliveryOfficeName"),
		Description: e.GetAttributeValue("description"),

		// Organization
		EmployeeID:     e.GetAttributeValue("employeeID"),
		EmployeeNumber: e.GetAttributeValue("employeeNumber"),
		Department:     e.GetAttributeValue("department"),
		Division:       e.GetAttributeValue("division"),
		Title:          e.GetAttributeValue("title"),
		Company:        e.GetAttributeValue("company"),
		Manager:        e.GetAttributeValue("manager"),

		// Address
		Street:     e.GetAttributeValue("streetAddress"),
		City:       e.GetAttributeValue("l"),
		State:      e.GetAttributeValue("st"),
		Country:    e.GetAttributeValue("co"),
		PostalCode: e.GetAttributeValue("postalCode"),

		// Directory
		ObjectGUID: e.GetAttributeValue("objectGUID"),
		ObjectSID:  e.GetAttributeValue("objectSid"),

		// Account
		AccountEnabled: (uac & 0x0002) == 0,
		AccountLocked: e.GetAttributeValue("lockoutTime") != "" &&
			e.GetAttributeValue("lockoutTime") != "0",

		// Membership
		Groups: e.GetAttributeValues("memberOf"),

		// Raw LDAP dates
		CreatedAt: ldapTime(e.GetAttributeValue("whenCreated")),
		UpdatedAt: ldapTime(e.GetAttributeValue("whenChanged")),
	}
}

func ldapTime(v string) time.Time {
	if v == "" {
		return time.Time{}
	}

	t, err := time.Parse("20060102150405.0Z", v)
	if err != nil {
		return time.Time{}
	}

	return t
}

// ----------------------------------------------------------------------
// Get User
// ----------------------------------------------------------------------

func GetUser(username string) (*config.User, error) {

	filter := fmt.Sprintf(
		"(|(sAMAccountName=%s)(userPrincipalName=%s))",
		gldap.EscapeFilter(username),
		gldap.EscapeFilter(username),
	)

	req := gldap.NewSearchRequest(
		config.LdapC.BaseDN,
		gldap.ScopeWholeSubtree,
		gldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		ldapAttributes,
		nil,
	)

	res, err := config.Conn.Search(req)
	if err != nil {
		return nil, err
	}

	if len(res.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := entryToUser(res.Entries[0])

	return &user, nil
}

// ----------------------------------------------------------------------
// Get All Users
// ----------------------------------------------------------------------

func GetAllUsers(conn *gldap.Conn, baseDN, ou string) ([]config.User, error) {

	searchBase := baseDN

	if strings.TrimSpace(ou) != "" {
		searchBase = fmt.Sprintf("OU=%s,%s", ou, baseDN)
	}

	req := gldap.NewSearchRequest(
		searchBase,
		gldap.ScopeWholeSubtree,
		gldap.NeverDerefAliases,
		0,
		0,
		false,
		"(&(objectCategory=person)(objectClass=user))",
		ldapAttributes,
		nil,
	)

	res, err := conn.Search(req)
	if err != nil {
		return nil, err
	}

	users := make([]config.User, 0, len(res.Entries))

	for _, e := range res.Entries {
		users = append(users, entryToUser(e))
	}

	return users, nil
}

// ----------------------------------------------------------------------
// Common Attributes
// ----------------------------------------------------------------------

var ldapAttributes = []string{
	"distinguishedName",
	"cn",
	"sAMAccountName",
	"userPrincipalName",
	"displayName",
	"givenName",
	"sn",
	"initials",

	"mail",
	"telephoneNumber",
	"mobile",
	"facsimileTelephoneNumber",
	"physicalDeliveryOfficeName",
	"description",

	"employeeID",
	"employeeNumber",
	"department",
	"division",
	"title",
	"company",
	"manager",

	"streetAddress",
	"l",
	"st",
	"co",
	"postalCode",

	"objectGUID",
	"objectSid",

	"userAccountControl",
	"lockoutTime",

	"memberOf",

	"whenCreated",
	"whenChanged",
}
