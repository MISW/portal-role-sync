package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mock_auth0 "github.com/MISW/portal-role-sync/infra/mock/auth0"
	"github.com/MISW/portal-role-sync/infra/portal"
)

func TestAuth0Reconciler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req := &ReconcileRequest{
		Members: portal.MemberRoles{
			"U111": {Role: "member"},
			"U222": {Role: "admin"},
			"U333": {Role: "retired"},
			"U444": {Role: "not_member"},
		},
	}

	ctrl := gomock.NewController(t)

	mockAuth0Client := mock_auth0.NewMockClient(ctrl)

	mockAuth0Client.EXPECT().UpdateRuleConfig(
		gomock.Any(),
		gomock.Eq("members"),
		JSONString(
			gomock.Eq(
				map[string]interface{}{
					"U111": "member",
					"U222": "admin",
				},
			),
		),
	).Return(nil)

	r := NewAuth0Reconciler(mockAuth0Client)

	if err := r.Reconcile(ctx, req); err != nil {
		t.Fatal(err)
	}
}

type jsonStringMatcher struct {
	m gomock.Matcher
}

type jsonStringFormatter struct {
	m gomock.Matcher
}

func JSONString(m gomock.Matcher) gomock.Matcher {
	return gomock.GotFormatterAdapter(&jsonStringFormatter{m}, &jsonStringMatcher{m})
}

func (j *jsonStringMatcher) Matches(x interface{}) bool {
	jsonStr, ok := x.(string)
	if !ok {
		return false
	}

	var value interface{}
	if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
		return false
	}

	return j.m.Matches(value)
}

func (j *jsonStringMatcher) String() string {
	return fmt.Sprintf("string(json.Marshal(%s))", j.m.String())
}

func (j *jsonStringFormatter) Got(got interface{}) string {
	var value interface{}
	json.Unmarshal([]byte(got.(string)), &value)

	f, ok := j.m.(gomock.GotFormatter)

	if ok {
		return fmt.Sprintf("string(json.Marshal(%s))", f.Got(value))
	}

	return fmt.Sprintf("string(json.Marshal(%#v))", value)
}
