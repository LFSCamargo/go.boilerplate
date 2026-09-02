package services

import "testing"

func TestStoreInterfacesAreSatisfiedByFakes(t *testing.T) {
	var users UserStore = newFakeUsers()
	var otps OTPStore = &fakeOTPs{}
	var revoked RevokedStore = newFakeRevoked()
	var mailer Mailer = NoopMailer{}
	if users == nil || otps == nil || revoked == nil || mailer == nil {
		t.Fatal("expected implementations")
	}
}
