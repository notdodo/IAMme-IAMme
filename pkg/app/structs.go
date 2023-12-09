package app

import "github.com/okta/okta-sdk-golang/v2/okta"

type Group struct {
	*okta.Group
	Members []*User
}

type User struct {
	*okta.User
}

type GroupRule struct {
	*okta.GroupRule
}
