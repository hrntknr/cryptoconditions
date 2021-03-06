package cryptoconditions

var (
	_ Fulfillment = FfPreimageSha256{}
	_ Fulfillment = new(FfPreimageSha256)

	_ Fulfillment                  = FfPrefixSha256{}
	_ Fulfillment                  = new(FfPrefixSha256)
	_ compoundConditionFulfillment = new(FfPrefixSha256)

	_ Fulfillment = FfEd25519Sha256{}
	_ Fulfillment = new(FfEd25519Sha256)

	_ Fulfillment = FfRsaSha256{}
	_ Fulfillment = new(FfRsaSha256)

	_ Fulfillment                  = FfThresholdSha256{}
	_ Fulfillment                  = new(FfThresholdSha256)
	_ compoundConditionFulfillment = new(FfThresholdSha256)
)
