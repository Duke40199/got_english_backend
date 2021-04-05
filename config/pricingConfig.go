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

type PricingNameConfig struct {
	MessagingSessionPricingID   string
	TranslationSessionPricingID string
	PrivateCallSessionPricingID string
}

func GetPricingNameConfig() *PricingNameConfig {
	return &PricingNameConfig{
		MessagingSessionPricingID:   "messaging",
		TranslationSessionPricingID: "translation",
		PrivateCallSessionPricingID: "private_call",
	}
}
