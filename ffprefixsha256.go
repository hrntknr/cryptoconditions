package cryptoconditions

import (
	"bytes"
	"crypto/sha256"

	"github.com/pkg/errors"
)

// FfPrefixSha256 implements the PREFIX-SHA-256 fulfillment.
type FfPrefixSha256 struct {
	Prefix           []byte `asn1:"tag:0"`
	MaxMessageLength uint32 `asn1:"tag:1"`

	// Only have either a sub-fulfillment or a sub-condition.
	SubFulfillment Fulfillment `asn1:"tag:2,choice:fulfillment"`
	subCondition   Condition   `asn1:"-"`
}

// NewPrefixSha256 creates a new PREFIX-SHA-256 fulfillment.
func NewPrefixSha256(prefix []byte, maxMessageLength uint32, subFf Fulfillment) *FfPrefixSha256 {
	return &FfPrefixSha256{
		Prefix:           prefix,
		MaxMessageLength: maxMessageLength,
		SubFulfillment:   subFf,
	}
}

// PrefixSha256Unfulfilled creates an unfulfilled PREFIX-SHA-256 fulfillment.
func NewPrefixSha256Unfulfilled(prefix []byte, maxMessageLength uint32, subCondition Condition) *FfPrefixSha256 {
	return &FfPrefixSha256{
		Prefix:           prefix,
		MaxMessageLength: maxMessageLength,
		subCondition:     subCondition,
	}
} //TODO consider if we really need this

func (f *FfPrefixSha256) ConditionType() ConditionType {
	return CTPrefixSha256
}

// SubCondition returns the sub-condition of this fulfillment.
func (f *FfPrefixSha256) SubCondition() Condition {
	if f.IsFulfilled() {
		return f.SubFulfillment.Condition()
	} else {
		return f.subCondition
	}
}

// IsFulfilled returns true if this fulfillment is fulfilled,
// i.e. when it contains a sub-fulfillment.
// If false, it only contains a sub-condition.
func (f *FfPrefixSha256) IsFulfilled() bool {
	return f.SubFulfillment != nil
}

func (f *FfPrefixSha256) fingerprintContents() []byte {
	content := struct {
		Prefix           []byte
		MaxMessageLength uint32
		SubCondition     Condition `asn1:"choice:condition"`
	}{
		Prefix:           f.Prefix,
		MaxMessageLength: f.MaxMessageLength,
		SubCondition:     f.SubCondition(),
	}

	encoded, err := ASN1Context.Encode(content)
	if err != nil {
		//TODO
		panic(err)
	}

	return encoded
}

func (f *FfPrefixSha256) fingerprint() []byte {
	hash := sha256.Sum256(f.fingerprintContents())
	return hash[:]
}

func (f *FfPrefixSha256) cost() int {
	encodedPrefix, err := ASN1Context.EncodeWithOptions(f.Prefix, "tag:0")
	if err != nil {
		//TODO
		panic(err)
	}
	return len(encodedPrefix) +
		int(f.MaxMessageLength) +
		f.SubCondition().Cost() +
		1024
}

func (f *FfPrefixSha256) subConditionTypeSet() ConditionTypeSet {
	var set ConditionTypeSet
	if f.IsFulfilled() {
		set = set.addRelevant(f.SubFulfillment)
	} else {
		set = set.addRelevant(f.subCondition)
	}
	return set
}

func (f *FfPrefixSha256) Condition() Condition {
	return NewCompoundCondition(f.ConditionType(), f.fingerprint(),
		f.cost(), f.subConditionTypeSet())
}

func (f *FfPrefixSha256) Encode() ([]byte, error) {
	return encodeFulfillment(f)
}

func (f *FfPrefixSha256) Validate(condition Condition, message []byte) error {
	if !matches(f, condition) {
		return fulfillmentDoesNotMatchConditionError
	}

	if !f.IsFulfilled() {
		return errors.New("cannot validate unfulfilled fulfillment.")
	}

	if len(message) > int(f.MaxMessageLength) {
		return errors.Errorf(
			"message length of %d exceeds limit of %d",
			len(message), f.MaxMessageLength)
	}

	buffer := new(bytes.Buffer)
	buffer.Write(f.Prefix)
	buffer.Write(message)
	newMessage := buffer.Bytes()

	return errors.Wrapf(f.SubFulfillment.Validate(nil, newMessage),
		"failed to validate sub-fulfillment with message %x", newMessage)
}
