
package password

import "testing"

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("SecurePassword123!")
	if err != nil {
		t.Fatal(err)
	}

	valid, err := Verify("SecurePassword123!", hash)
	if err != nil || !valid {
		t.Fatal("correct password should verify")
	}

	valid, err = Verify("WrongPassword", hash)
	if err != nil || valid {
		t.Fatal("wrong password should not verify")
	}
}
