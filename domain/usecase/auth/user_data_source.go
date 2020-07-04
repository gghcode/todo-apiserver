// UserDataSource retrieves user data from detail source
type UserDataSource interface {
	UserByUserName(username string) (interface{}, error)
}