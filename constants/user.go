package constants

const (
	CtxAuthenticatedUserKey = "CtxAuthenticatedUserKey"
	Admin                   = "admin"
	User                    = "user"
	Male                    = "male"
	Female                  = "female"
)

var (
	MapperGenderToId = map[string]int{
		Male:   1,
		Female: 2,
	}

	MapperRoleToId = map[string]int{
		Admin: 1,
		User:  2,
	}

	ListGender = []string{Male, Female}
)
