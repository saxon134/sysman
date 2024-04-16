package rbac

type User struct {
	Id    int64  `json:"id"`
	Roles string `json:"roles"`
}
