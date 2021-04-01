package config

type PricingIDConfig struct {
	MessagingSessionPricingID   uint
	TranslationSessionPricingID uint
	PrivateCallSessionPricingID uint
}

func GetPricingIDConfig() *PricingIDConfig {
	return &PricingIDConfig{
		MessagingSessionPricingID:   1,
		TranslationSessionPricingID: 2,
		PrivateCallSessionPricingID: 3,
	}
}
