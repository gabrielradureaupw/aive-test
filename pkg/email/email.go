package email

import "github.com/rs/zerolog/log"

func Send(email, subject, msg string) error {
	log.Info().Msgf(`Sending %s to %s with content: '%s'`, subject, email, msg)
	return nil
}
