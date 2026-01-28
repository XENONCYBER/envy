package tests

import (
	"bytes"
	"testing"

	"envy/internal/crypto"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}

	if len(salt1) != 16 { // saltSize is 16
		t.Errorf("GenerateSalt() returned %d bytes, want 16", len(salt1))
	}

	// Generate another salt to ensure randomness
	salt2, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() second call error: %v", err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Error("GenerateSalt() returned identical salts, expected unique values")
	}
}

func TestDeriveKey(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	password := "testpassword123"

	key1 := crypto.DeriveKey(password, salt)

	if len(key1) != 32 { // keySize is 32
		t.Errorf("DeriveKey() returned %d bytes, want 32", len(key1))
	}

	// Same password and salt should produce same key
	key2 := crypto.DeriveKey(password, salt)
	if !bytes.Equal(key1, key2) {
		t.Error("DeriveKey() should produce identical keys for same password and salt")
	}

	// Different password should produce different key
	key3 := crypto.DeriveKey("differentpassword", salt)
	if bytes.Equal(key1, key3) {
		t.Error("DeriveKey() should produce different keys for different passwords")
	}

	// Different salt should produce different key
	salt2, _ := crypto.GenerateSalt()
	key4 := crypto.DeriveKey(password, salt2)
	if bytes.Equal(key1, key4) {
		t.Error("DeriveKey() should produce different keys for different salts")
	}
}

func TestGenerateAuthHash(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes

	hash1 := crypto.GenerateAuthHash(key)
	hash2 := crypto.GenerateAuthHash(key)

	if hash1 != hash2 {
		t.Error("GenerateAuthHash() should produce identical hashes for same key")
	}

	// Different key should produce different hash
	differentKey := []byte("abcdef0123456789abcdef0123456789")
	hash3 := crypto.GenerateAuthHash(differentKey)

	if hash1 == hash3 {
		t.Error("GenerateAuthHash() should produce different hashes for different keys")
	}
}

func TestVerifyAuthHash(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	hash := crypto.GenerateAuthHash(key)

	if !crypto.VerifyAuthHash(key, hash) {
		t.Error("VerifyAuthHash() should return true for matching key and hash")
	}

	wrongKey := []byte("wrongkey0123456789abcdef01234567")
	if crypto.VerifyAuthHash(wrongKey, hash) {
		t.Error("VerifyAuthHash() should return false for wrong key")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes for AES-256
	plaintext := []byte("Hello, World! This is a secret message.")

	// Encrypt
	ciphertext, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	if ciphertext == "" {
		t.Error("Encrypt() returned empty ciphertext")
	}

	// Decrypt
	decrypted, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	plaintext := []byte("Same message")

	// Encrypt the same message twice
	ciphertext1, _ := crypto.Encrypt(plaintext, key)
	ciphertext2, _ := crypto.Encrypt(plaintext, key)

	// Should produce different ciphertexts due to random nonce
	if ciphertext1 == ciphertext2 {
		t.Error("Encrypt() should produce different ciphertexts for same plaintext (random nonce)")
	}

	// But both should decrypt to the same plaintext
	decrypted1, _ := crypto.Decrypt(ciphertext1, key)
	decrypted2, _ := crypto.Decrypt(ciphertext2, key)

	if !bytes.Equal(decrypted1, decrypted2) {
		t.Error("Both ciphertexts should decrypt to same plaintext")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	wrongKey := []byte("wrongkey0123456789abcdef01234567")
	plaintext := []byte("Secret message")

	ciphertext, _ := crypto.Encrypt(plaintext, key)

	_, err := crypto.Decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Error("Decrypt() should return error with wrong key")
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")

	// Invalid base64
	_, err := crypto.Decrypt("not-valid-base64!!!", key)
	if err == nil {
		t.Error("Decrypt() should return error for invalid base64")
	}

	// Too short ciphertext
	_, err = crypto.Decrypt("YWJj", key) // "abc" in base64
	if err == nil {
		t.Error("Decrypt() should return error for too short ciphertext")
	}
}

func TestEncryptEmptyPlaintext(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	plaintext := []byte("")

	ciphertext, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error with empty plaintext: %v", err)
	}

	decrypted, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt() error with empty plaintext: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("Decrypt() = %q, want empty", decrypted)
	}
}

func TestEncryptLargePlaintext(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	plaintext := make([]byte, 1024*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	ciphertext, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error with large plaintext: %v", err)
	}

	decrypted, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt() error with large plaintext: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypt() did not return original large plaintext")
	}
}
