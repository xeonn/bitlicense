package bitlicense

import (
	"fmt"
	"testing"
	"time"

	"gitlab.com/bitify-pub/byutils/timeutils"
)

func TestIssue(t *testing.T) {
	expiryTime, err := time.Parse(time.RFC3339, "2021-12-31T23:59:59Z")
	if err != nil {
		fmt.Println(err)
		return
	}

	// expiry time must be rounded to the previous 5 minutes
	rounded := timeutils.RoundDownTo5Minutes(expiryTime.UTC())

	lic1 := issue("demo", rounded.Format("2006-01-02"), "../certs/privatekey")

	lic2 := issue("demo", rounded.Add(48 * time.Hour).Format("2006-01-02"), "../certs/privatekey")

	if lic1 == nil {
		t.Fail()
	}

	if lic2 == nil {
		t.Fail()
	} else if lic1.Signature != lic2.Signature {
		t.Log("generation successful")
		t.Log("Encoded Signature (base64): ", rounded, lic1.Signature)
		t.Log("Encoded Signature2 (base64): ", rounded.Add(5* time.Hour), lic2.Signature)
	} else {
		t.Errorf("generation failed")
	}
}