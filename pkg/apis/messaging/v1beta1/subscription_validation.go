/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"

	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/api/equality"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmp"
)

func (s *Subscription) Validate(ctx context.Context) *apis.FieldError {
	withNS := apis.WithinParent(ctx, s.ObjectMeta)
	return s.Spec.Validate(withNS).ViaField("spec")
}

func (spec *SubscriptionSpec) Validate(ctx context.Context) *apis.FieldError {
	// We require always Channel.
	// Also at least one of 'subscriber' and 'reply' must be defined (non-nil and non-empty).

	var errs *apis.FieldError

	// Validate the Channel
	if spec.Channel != nil {
		if ce := spec.Channel.Validate(ctx); ce != nil {
			errs = errs.Also(ce.ViaField("channel"))
		}
	}

	missingSubscriber := isDestinationNilOrEmpty(spec.Subscriber)
	missingReply := isDestinationNilOrEmpty(spec.Reply)
	if missingSubscriber && missingReply {
		fe := apis.ErrMissingField("reply", "subscriber")
		fe.Details = "the Subscription must reference at least one of (reply or a subscriber)"
		errs = errs.Also(fe)
	}

	if !missingSubscriber {
		if fe := spec.Subscriber.Validate(ctx); fe != nil {
			errs = errs.Also(fe.ViaField("subscriber"))
		}
	}

	if !missingReply {
		if fe := spec.Reply.Validate(ctx); fe != nil {
			errs = errs.Also(fe.ViaField("reply"))
		}
	}

	return errs
}

func isDestinationNilOrEmpty(d *duckv1.Destination) bool {
	return d == nil || equality.Semantic.DeepEqual(d, &duckv1.Destination{})
}

func (s *Subscription) CheckImmutableFields(ctx context.Context, original *Subscription) *apis.FieldError {
	if original == nil {
		return nil
	}

	// Only Subscriber and Reply are mutable.
	ignoreArguments := cmpopts.IgnoreFields(SubscriptionSpec{}, "Subscriber", "Reply")
	if diff, err := kmp.ShortDiff(original.Spec, s.Spec, ignoreArguments); err != nil {
		return &apis.FieldError{
			Message: "Failed to diff Subscription",
			Paths:   []string{"spec"},
			Details: err.Error(),
		}
	} else if diff != "" {
		return &apis.FieldError{
			Message: "Immutable fields changed (-old +new)",
			Paths:   []string{"spec"},
			Details: diff,
		}
	}
	return nil
}
