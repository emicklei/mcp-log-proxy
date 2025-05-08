package main

type proxyInstance struct {
	host    string
	port    int
	title   string
	command string
}

func register(inst proxyInstance) error {
	return nil
}
