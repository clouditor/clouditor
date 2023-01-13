package testdata

import "clouditor.io/clouditor/api/auth"

func NewUser1() *auth.User {
	return &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}
}

func NewUser2() *auth.User {
	return &auth.User{
		Username: "SomeName2",
		Password: "SomePassword2",
		Email:    "SomeMail2",
		FullName: "SomeFullName2",
	}
}
