package app

import "github.com/okta/okta-sdk-golang/v2/okta"

type Group struct {
	Id string
	*okta.Group
	Members []*User
}

type User struct {
	Id string
	*okta.User
}

type GroupRule struct {
	Id string
	*okta.GroupRule
}
