package services

type UserService interface {
	OnlineCheckByUserId(userId string) bool

	OnlineCheckByUsername(username string) bool

	Logout(userId string) bool
}
