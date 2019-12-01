package schema

type Human struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Dogs []*Dog `json:"dogs"`
}
