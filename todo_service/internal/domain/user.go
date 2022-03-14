package domain

//User is a user of system. The Todos array is only stored undone todo.
type User struct {
	Email    string `bson:"_id,omitempty"`
	Name     string `bson:"name,omitempty"`
	Password string `bson:"password,omitempty"`
	Todos    []Todo `bson:"todos"`
	IsActive bool   `bson:"is_active"`
}

//NewUser creates new user.
func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// AddNewTodo add new todo to the todo collection.
func (u *User) AddNewTodo(todo Todo) {
	u.Todos = append(u.Todos, todo)
}
