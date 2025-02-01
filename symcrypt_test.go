package symcrypt_test

import (
	"testing"

	"github.com/peterldowns/testy/assert"
	"github.com/peterldowns/testy/check"

	"github.com/peterldowns/symcrypt"
)

func TestRoundtripSucceedsForSameOwner(t *testing.T) {
	t.Parallel()
	secretKey, err := symcrypt.GenerateRandomKey()
	assert.Nil(t, err)

	symc, err := symcrypt.NewClient(secretKey)
	assert.Nil(t, err)

	plaintext := symcrypt.Plaintext("ascx_mysecretaccesstoken")
	userID := symcrypt.Owner("userid_000111")

	encrypted, err := symc.Encrypt(plaintext, userID)
	assert.Nil(t, err)

	decrypted, err := symc.Decrypt(encrypted, userID)
	assert.Nil(t, err)

	assert.Equal(t, plaintext, decrypted)
}

func TestAssociatedDataMustMatch(t *testing.T) {
	t.Parallel()
	secretKey, err := symcrypt.GenerateRandomKey()
	assert.Nil(t, err)

	symc, err := symcrypt.NewClient(secretKey)
	assert.Nil(t, err)

	plaintext := symcrypt.Plaintext("ascx_mysecretaccesstoken")

	// A value encrypted for one user should fail to decrypt for another user.
	userA := symcrypt.Owner("userid_000111")
	encrypted, err := symc.Encrypt(plaintext, userA)
	assert.Nil(t, err)

	userB := symcrypt.Owner("userid_222333")
	decrypted, err := symc.Decrypt(encrypted, userB)
	check.Error(t, err)
	check.NotEqual(t, plaintext, decrypted)
	check.Equal(t, "", string(decrypted))

	emptyUser := symcrypt.Owner("")
	decrypted, err = symc.Decrypt(encrypted, emptyUser)
	check.Error(t, err)
	check.NotEqual(t, plaintext, decrypted)
	check.Equal(t, "", string(decrypted))
}

func TestEncryptingEmptyPlaintextSucceeds(t *testing.T) {
	t.Parallel()
	secretKey, err := symcrypt.GenerateRandomKey()
	assert.Nil(t, err)

	symc, err := symcrypt.NewClient(secretKey)
	assert.Nil(t, err)

	plaintext := symcrypt.Plaintext("")
	owner := symcrypt.Owner("user_001")

	encrypted, err := symc.Encrypt(plaintext, owner)
	assert.Nil(t, err)
	decrypted, err := symc.Decrypt(encrypted, owner)
	assert.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptingWithEmptyOwnerSucceeds(t *testing.T) {
	t.Parallel()
	secretKey, err := symcrypt.GenerateRandomKey()
	assert.Nil(t, err)

	symc, err := symcrypt.NewClient(secretKey)
	assert.Nil(t, err)

	plaintext := symcrypt.Plaintext("my super secret")
	owner := symcrypt.Owner("")

	encrypted, err := symc.Encrypt(plaintext, owner)
	assert.Nil(t, err)
	decrypted, err := symc.Decrypt(encrypted, owner)
	assert.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}

/*
// Uncomment this test to show that using plain strings
// will cause errors at compile time.
func TestThatStringsDontWork(t *testing.T) {

	t.Parallel()
	secretKey, err := symcrypt.GenerateRandomKey()
	assert.Nil(t, err)

	symc, err := symcrypt.NewClient(secretKey)
	assert.Nil(t, err)

	var plaintext string = "my super secret"
	var owner string = "some owner ID"

	encrypted, err := symc.Encrypt(plaintext, owner)
	assert.Nil(t, err)
	decrypted, err := symc.Decrypt(encrypted, owner)
	assert.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}
*/

/*
// Uncomment this test to show that mis-matching the Owner/Plaintext/Ciphertext
// types will cause errors at compile time.
func TestThatStringsDontWork(t *testing.T) {
	t.Parallel()
	secretKey, err := symcrypt.GenerateRandomKey()
	assert.Nil(t, err)

	symc, err := symcrypt.NewClient(secretKey)
	assert.Nil(t, err)

	// should be Plaintext, not Ciphertext
	var plaintext symcrypt.Ciphertext = "my super secret"
	// should be Owner, not Plaintext
	var owner symcrypt.Plaintext = "some owner ID"

	encrypted, err := symc.Encrypt(plaintext, owner)
	assert.Nil(t, err)
	decrypted, err := symc.Decrypt(encrypted, owner)
	assert.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}
*/
