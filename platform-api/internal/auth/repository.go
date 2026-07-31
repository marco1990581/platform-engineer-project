


type UserRepository interface {
	GetByUserName(username string) (*models.User,error)

}
