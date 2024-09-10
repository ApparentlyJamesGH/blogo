package nostr

import (
	"fmt"

	gonostr "github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// Returns a key set in the following order: sk, pk, nsec, npub
func GetNewKeySet() (string, string, string, string, error) {
	sk := gonostr.GeneratePrivateKey()
	pk, err := gonostr.GetPublicKey(sk)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to get public key: %w", err)
	}
	nsec, err := nip19.EncodePrivateKey(sk)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to encode public key: %w", err)
	}
	npub, err := nip19.EncodePublicKey(pk)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to encode public key: %w", err)
	}
	return sk, pk, nsec, npub, nil
}
