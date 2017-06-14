package cryptoconditions

import (
	"reflect"

	"github.com/pkg/errors"
)

// Fulfillment defines the fulfillment interface.
type Fulfillment interface {
	// ConditionType returns the type of condition this fulfillment fulfills.
	ConditionType() ConditionType

	// Condition generates the condition that this fulfillment fulfills.
	Condition() Condition
	//TODO consider moving Condition away from here because the next two can make up for it

	// Encode encodes the fulfillment into binary format.
	Encode() ([]byte, error)

	// Validate checks whether this fulfillment correctly validates the given
	// condition using the specified message.
	// It returns nil if it does, an error indicating the problem otherwise.
	Validate(Condition, []byte) error

	// fingerprint calculates the fingerprint of the condition this fulfillment
	// fulfills.
	fingerprint() []byte

	// cost calculates the cost metric of this fulfillment.
	cost() int
}

// compoundConditionFulfillment is an interface that fulfillments for compound
// conditions have to implement to be able to indicate the condition types of
// their sub-fulfillments.
type compoundConditionFulfillment interface {
	// subConditionsTypeSet returns the set with all the different types
	// amongst sub-conditions of this fulfillment.
	subConditionsTypeSet() ConditionTypeSet
}

// fulfillmentTypeMap maps ConditionTypes to the corresponding Go type for the
// fulfillment for that condition.
var fulfillmentTypeMap = map[ConditionType]reflect.Type{
	CTEd25519Sha256:   reflect.TypeOf(FfEd25519Sha256{}),
	CTPrefixSha256:    reflect.TypeOf(FfPrefixSha256{}),
	CTPreimageSha256:  reflect.TypeOf(FfPreimageSha256{}),
	CTRsaSha256:       reflect.TypeOf(FfRsaSha256{}),
	CTThresholdSha256: reflect.TypeOf(FfThresholdSha256{}),
}

// fulfillmentDoesNotMatchConditionError is the error we throw when trying to
// validate a condition with a wrong
// fulfillment.
var fulfillmentDoesNotMatchConditionError = errors.New(
	"The fulfillment does not match the given condition")
