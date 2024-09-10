package nostr

import (
	"fmt"

	gonostr "github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Init() error {
	nsec := viper.GetString("nostr.nsec")
	var npub string
	var err error
	if nsec == "" {
		log.Warn().Msg("NOSTR_NSEC not set. Generating a new key pair.")
		nostrSk, nostrPk, nsec, npub, err = GetNewKeySet()
		if err != nil {
			return fmt.Errorf("failed to get public key: %w", err)
		}
		log.Info().Msgf("Generated new key set:\n\nnsec: %v\nnpub: %v\n\n", nsec, npub)
	} else {
		log.Info().Msg("NOSTR_NSEC set. Deriving existing key pair.")
		_, value, err := nip19.Decode(nsec)
		if err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}

		var ok bool
		nostrSk, ok = value.(string)
		if !ok {
			return fmt.Errorf("failed to nip19 decode to string: %w", err)
		}

		nostrPk, err = gonostr.GetPublicKey(nostrSk)
		if err != nil {
			return fmt.Errorf("failed to get public key: %w", err)
		}
		npub, err = nip19.EncodePublicKey(nostrPk)
		if err != nil {
			return fmt.Errorf("failed to encode public key: %w", err)
		}
	}

	fmt.Println("Public Key:", nostrPk)
	fmt.Println("npub: ", npub)
	return nil
}
