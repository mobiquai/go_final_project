package taskscontrol

const dateLayout string = "20060102"

type TaskIdResponse struct {
	Id int `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SignIn struct {
	Password string `json:"password"`
}

type AuthToken struct {
	Token string `json:"token"`
}
