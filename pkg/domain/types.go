package domain

type UserId string

func (userId *UserId) String() string {
	return string(*userId)
}

type Email string

func (email *Email) String() string {
	return string(*email)
}

type Password string

func (password *Password) String() string {
	return string(*password)
}

type Role string

func (role *Role) String() string {
	return string(*role)
}

type FullName string

func (fullName *FullName) String() string {
	return string(*fullName)
}
