package handlers

type Result struct {
	Payload []byte
	Status  int
}

type TextResult struct {
	Payload string
	Status  int
}
