package enum

type UserStatus int8

const (
	INACTIVE UserStatus = 1
	ACTIVE   UserStatus = 2
	SUSPEND  UserStatus = 3
)

var strings []string = []string{"", "INACTIVE", "ACTIVE", "SUSPEND"}

func (u UserStatus) String() string {
	return strings[u]
}
