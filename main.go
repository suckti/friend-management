package main

type User struct {
	Email     string   `json:"email"`
	Friends   []string `json:"friends"`
	Block     []string `json:"block"`
	Subscribe []string `json:"subscribe"`
}

type Friend struct {
	Friends []string
}

type Subscribe struct {
	Requestor string
	Target    string
}

type Notification struct {
	Sender string
	Text   string
}

type Person struct {
	Name  string
	Phone string
}

func main() {

}
