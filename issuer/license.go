package issuer

type License struct {
	Client    string
	Expiry    string
	Signature string // base64 encoded string
}
