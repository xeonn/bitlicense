package issuer

import "testing"

func TestValidate(t *testing.T) {
	lic := &License{
		Client: "demo",
		Expiry: "2024-12-06T00:00:00Z",
		Signature: "fIEU4wgnHwR1YPOygmMWg2BJmnZnaVYc4WBl0oKrfdESBPUocVLH+XJNUwt9CdxjjBD5X6giUVB2ITM//OKUDg==",
	}

	if !Validate(lic) {
		t.Fail()
	}
}