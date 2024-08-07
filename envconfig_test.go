// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/reMarkable/envconfig/v2/types"
)

type HonorDecodeInStruct struct {
	Value string
}

func (h *HonorDecodeInStruct) Decode(env string) error {
	h.Value = "decoded"
	return nil
}

type CustomURL struct {
	Value *url.URL
}

func (cu *CustomURL) UnmarshalBinary(data []byte) error {
	u, err := url.Parse(string(data))
	cu.Value = u
	return err
}

type Specification struct {
	Embedded                     `envconfig:"EMBEDDED" desc:"can we document a struct"`
	EmbeddedButIgnored           `envconfig:"EMBEDDED_BUT_IGNORED" ignored:"true"`
	Debug                        bool           `envconfig:"DEBUG"`
	Port                         int            `envconfig:"PORT"`
	Rate                         float32        `envconfig:"RATE"`
	User                         string         `envconfig:"USER"`
	TTL                          uint32         `envconfig:"TTL"`
	Timeout                      time.Duration  `envconfig:"TIMEOUT"`
	AdminUsers                   []string       `envconfig:"ADMINUSERS"`
	MagicNumbers                 []int          `envconfig:"MAGICNUMBERS"`
	EmptyNumbers                 []int          `envconfig:"EMPTYNUMBERS"`
	ByteSlice                    []byte         `envconfig:"BYTESLICE"`
	ColorCodes                   map[string]int `envconfig:"COLORCODES"`
	SomePointer                  *string        `envconfig:"SOMEPOINTER"`
	SomePointerWithDefault       *string        `envconfig:"SOMEPOINTERWITHDEFAULT" default:"foo2baz" desc:"foorbar is the word"`
	MultiWordVarWithAlt          string         `envconfig:"MULTI_WORD_VAR_WITH_ALT" desc:"what alt"`
	MultiWordVarWithLowerCaseAlt string         `envconfig:"multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt              string         `envconfig:"SERVICE_HOST"`
	DefaultVar                   string         `envconfig:"DEFAULTVAR" default:"foobar"`
	RequiredVar                  string         `envconfig:"REQUIREDVAR" required:"True"`
	NoPrefixDefault              string         `envconfig:"BROKER" default:"127.0.0.1"`
	RequiredDefault              string         `envconfig:"REQUIREDDEFAULT" required:"true" default:"foo2bar"`
	Ignored                      string         `envconfig:"IGNORED" ignored:"true"`
	NestedSpecification          struct {
		Property            string `envconfig:"inner"`
		PropertyWithDefault string `envconfig:"PROPERTYWITHDEFAULT" default:"fuzzybydefault"`
	} `envconfig:"outer"`
	AfterNested                    string                        `envconfig:"AFTERNESTED"`
	DecodeStruct                   HonorDecodeInStruct           `envconfig:"honor"`
	Datetime                       time.Time                     `envconfig:"DATETIME"`
	MapField                       map[string]string             `envconfig:"MAPFIELD" default:"one:two;three:four"`
	EmptyMapField                  map[string]string             `envconfig:"EMPTY_MAPFIELD"`
	UrlValue                       CustomURL                     `envconfig:"URLVALUE"`
	UrlPointer                     *CustomURL                    `envconfig:"URLPOINTER"`
	GooglePubSubTopic              types.GooglePubSubTopic       `envconfig:"GOOGLE_PUBSUB_TOPIC"`
	GoogleFirestoreDatabase        types.GoogleFirestoreDatabase `envconfig:"GOOGLE_FIRESTORE_DATABASE"`
	GoogleFirestoreDatabaseDefault types.GoogleFirestoreDatabase `envconfig:"GOOGLE_FIRESTORE_DATABASE_DEFAULT"`
}

type Embedded struct {
	Enabled             bool   `envconfig:"ENABLED" desc:"some embedded value"`
	EmbeddedPort        int    `envconfig:"EMBEDDEDPORT"`
	MultiWordVar        string `envconfig:"MULTIWORDVAR"`
	MultiWordVarWithAlt string `envconfig:"MULTI_WITH_DIFFERENT_ALT"`
	EmbeddedAlt         string `envconfig:"EMBEDDED_WITH_ALT"`
	EmbeddedIgnored     string `envconfig:"EMBEDDEDIGNORED" ignored:"true"`
}

type EmbeddedButIgnored struct {
	FirstEmbeddedButIgnored  string `envconfig:"FIRST_EMBEDDED_BUT_IGNORED"`
	SecondEmbeddedButIgnored string `envconfig:"SECOND_EMBEDDED_BUT_IGNORED"`
}

func TestProcess(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("ENV_CONFIG_TIMEOUT", "2m")
	os.Setenv("ENV_CONFIG_ADMINUSERS", "John,Adam,Will")
	os.Setenv("ENV_CONFIG_MAGICNUMBERS", "5,10,20")
	os.Setenv("ENV_CONFIG_EMPTYNUMBERS", "")
	os.Setenv("ENV_CONFIG_BYTESLICE", "dGhpcyBpcyBhIHRlc3QgdmFsdWU=")
	os.Setenv("ENV_CONFIG_COLORCODES", "red:1;green:2;blue:3")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_TTL", "30")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	os.Setenv("ENV_CONFIG_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_OUTER_INNER", "iamnested")
	os.Setenv("ENV_CONFIG_AFTERNESTED", "after")
	os.Setenv("ENV_CONFIG_HONOR", "honor")
	os.Setenv("ENV_CONFIG_DATETIME", "2016-08-16T18:57:05Z")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_AUTO_SPLIT", "24")
	os.Setenv("ENV_CONFIG_MULTI_WORD_ACR_WITH_AUTO_SPLIT", "25")
	os.Setenv("ENV_CONFIG_URLVALUE", "https://github.com/kelseyhightower/envconfig")
	os.Setenv("ENV_CONFIG_URLPOINTER", "https://github.com/kelseyhightower/envconfig")
	os.Setenv("ENV_CONFIG_GOOGLE_PUBSUB_TOPIC", "projects/project-id/topics/topic-id")
	os.Setenv("ENV_CONFIG_GOOGLE_FIRESTORE_DATABASE", "projects/project-id/databases/db")
	os.Setenv("ENV_CONFIG_GOOGLE_FIRESTORE_DATABASE_DEFAULT", "projects/project-id/databases/(default)")
	err := Process("env_config", &s)
	if err != nil {
		t.Error(err.Error())
	}
	// This is an inversion of the original test, since we have removed the
	// fallback of the Alt keyword, it no longer magically reads the non-prefixed
	// version.
	if s.NoPrefixWithAlt != "" {
		t.Errorf("expected %v, got %v", "", s.NoPrefixWithAlt)
	}
	if !s.Debug {
		t.Errorf("expected %v, got %v", true, s.Debug)
	}
	if s.Port != 8080 {
		t.Errorf("expected %d, got %v", 8080, s.Port)
	}
	if s.Rate != 0.5 {
		t.Errorf("expected %f, got %v", 0.5, s.Rate)
	}
	if s.TTL != 30 {
		t.Errorf("expected %d, got %v", 30, s.TTL)
	}
	if s.User != "Kelsey" {
		t.Errorf("expected %s, got %s", "Kelsey", s.User)
	}
	if s.Timeout != 2*time.Minute {
		t.Errorf("expected %s, got %s", 2*time.Minute, s.Timeout)
	}
	if s.RequiredVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.RequiredVar)
	}
	if len(s.AdminUsers) != 3 ||
		s.AdminUsers[0] != "John" ||
		s.AdminUsers[1] != "Adam" ||
		s.AdminUsers[2] != "Will" {
		t.Errorf("expected %#v, got %#v", []string{"John", "Adam", "Will"}, s.AdminUsers)
	}
	if len(s.MagicNumbers) != 3 ||
		s.MagicNumbers[0] != 5 ||
		s.MagicNumbers[1] != 10 ||
		s.MagicNumbers[2] != 20 {
		t.Errorf("expected %#v, got %#v", []int{5, 10, 20}, s.MagicNumbers)
	}
	if len(s.EmptyNumbers) != 0 {
		t.Errorf("expected %#v, got %#v", []int{}, s.EmptyNumbers)
	}
	expected := "this is a test value"
	if string(s.ByteSlice) != expected {
		t.Errorf("expected %v, got %v", expected, string(s.ByteSlice))
	}
	if s.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}

	if len(s.ColorCodes) != 3 ||
		s.ColorCodes["red"] != 1 ||
		s.ColorCodes["green"] != 2 ||
		s.ColorCodes["blue"] != 3 {
		t.Errorf(
			"expected %#v, got %#v",
			map[string]int{
				"red":   1,
				"green": 2,
				"blue":  3,
			},
			s.ColorCodes,
		)
	}

	if s.NestedSpecification.Property != "iamnested" {
		t.Errorf("expected '%s' string, got %#v", "iamnested", s.NestedSpecification.Property)
	}

	if s.NestedSpecification.PropertyWithDefault != "fuzzybydefault" {
		t.Errorf("expected default '%s' string, got %#v", "fuzzybydefault", s.NestedSpecification.PropertyWithDefault)
	}

	if s.AfterNested != "after" {
		t.Errorf("expected default '%s' string, got %#v", "after", s.AfterNested)
	}

	if s.DecodeStruct.Value != "decoded" {
		t.Errorf("expected default '%s' string, got %#v", "decoded", s.DecodeStruct.Value)
	}

	if expected := time.Date(2016, 8, 16, 18, 57, 05, 0, time.UTC); !s.Datetime.Equal(expected) {
		t.Errorf("expected %s, got %s", expected.Format(time.RFC3339), s.Datetime.Format(time.RFC3339))
	}

	u, err := url.Parse("https://github.com/kelseyhightower/envconfig")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if *s.UrlValue.Value != *u {
		t.Errorf("expected %q, got %q", u, s.UrlValue.Value.String())
	}

	if *s.UrlPointer.Value != *u {
		t.Errorf("expected %q, got %q", u, s.UrlPointer.Value.String())
	}

	if s.GooglePubSubTopic.ProjectID != "project-id" {
		t.Errorf("expected %s, got %s", "project-id", s.GooglePubSubTopic.ProjectID)
	}

	if s.GooglePubSubTopic.TopicID != "topic-id" {
		t.Errorf("expected %s, got %s", "topic-id", s.GooglePubSubTopic.TopicID)
	}

	if s.GoogleFirestoreDatabase.ProjectID != "project-id" {
		t.Errorf("expected %s, got %s", "project-id", s.GoogleFirestoreDatabase.ProjectID)
	}

	if s.GoogleFirestoreDatabase.Database != "db" {
		t.Errorf("expected %s, got %s", "db", s.GoogleFirestoreDatabase.Database)
	}

	if s.GoogleFirestoreDatabaseDefault.ProjectID != "project-id" {
		t.Errorf("expected %s, got %s", "project-id", s.GoogleFirestoreDatabaseDefault.ProjectID)
	}

	if s.GoogleFirestoreDatabaseDefault.Database != "(default)" {
		t.Errorf("expected %s, got %s", "default", s.GoogleFirestoreDatabaseDefault.Database)
	}
}

func TestParseErrorBool(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "string")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Debug" {
		t.Errorf("expected %s, got %v", "Debug", v.FieldName)
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}
}

func TestParseErrorFloat32(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_RATE", "string")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Rate" {
		t.Errorf("expected %s, got %v", "Rate", v.FieldName)
	}
	if s.Rate != 0 {
		t.Errorf("expected %v, got %v", 0, s.Rate)
	}
}

func TestParseErrorInt(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_PORT", "string")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Port" {
		t.Errorf("expected %s, got %v", "Port", v.FieldName)
	}
	if s.Port != 0 {
		t.Errorf("expected %v, got %v", 0, s.Port)
	}
}

func TestParseErrorUint(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_TTL", "-30")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "TTL" {
		t.Errorf("expected %s, got %v", "TTL", v.FieldName)
	}
	if s.TTL != 0 {
		t.Errorf("expected %v, got %v", 0, s.TTL)
	}
}

func TestParseErrorGooglePubSubTopic(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_GOOGLE_PUBSUB_TOPIC", "invalid/project-id/topics")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}

	if v.FieldName != "GooglePubSubTopic" {
		t.Errorf("expected %s, got %v", "GooglePubSubTopic", v.FieldName)
	}

	if s.GooglePubSubTopic.TopicID != "" {
		t.Errorf("expected %s, got %s", "", s.GooglePubSubTopic.TopicID)
	}

	if s.GooglePubSubTopic.ProjectID != "" {
		t.Errorf("expected %s, got %s", "", s.GooglePubSubTopic.ProjectID)
	}

	if v.Err != types.ErrInvalidGoogleTopicID {
		t.Errorf("unexpected %s, got %s", types.ErrInvalidGoogleTopicID, v.Err)
	}
}

func TestParseErrorGoogleFirestoreDatabase(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_GOOGLE_FIRESTORE_DATABASE", "invalid/project-id/databases")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}

	if v.FieldName != "GoogleFirestoreDatabase" {
		t.Errorf("expected %s, got %v", "GoogleFirestoreDatabase", v.FieldName)
	}

	if s.GoogleFirestoreDatabase.Database != "" {
		t.Errorf("expected %s, got %s", "", s.GoogleFirestoreDatabase.Database)
	}

	if s.GoogleFirestoreDatabase.ProjectID != "" {
		t.Errorf("expected %s, got %s", "", s.GoogleFirestoreDatabase.ProjectID)
	}

	if v.Err != types.ErrInvalidGoogleFirestoreID {
		t.Errorf("unexpected %s, got %s", types.ErrInvalidGoogleFirestoreID, v.Err)
	}
}

func TestErrInvalidSpecification(t *testing.T) {
	m := make(map[string]string)
	err := Process("env_config", &m)
	if err != ErrInvalidSpecification {
		t.Errorf("expected %v, got %v", ErrInvalidSpecification, err)
	}
}

func TestUnsetVars(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("USER", "foo")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// If the var is not defined the non-prefixed version should not be used
	// unless the struct tag says so
	if s.User != "" {
		t.Errorf("expected %q, got %q", "", s.User)
	}
}

func TestAlternateVarNames(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR", "foo")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_LOWER_CASE_ALT", "baz")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// Setting the alt version of the var in the environment has no effect if
	// the struct tag is not supplied
	if s.MultiWordVar != "" {
		t.Errorf("expected %q, got %q", "", s.MultiWordVar)
	}

	// Setting the alt version of the var in the environment correctly sets
	// the value if the struct tag IS supplied
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %q, got %q", "bar", s.MultiWordVarWithAlt)
	}

	// Alt value is not case sensitive and is treated as all uppercase
	if s.MultiWordVarWithLowerCaseAlt != "baz" {
		t.Errorf("expected %q, got %q", "baz", s.MultiWordVarWithLowerCaseAlt)
	}
}

func TestRequiredVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foobar")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.RequiredVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.RequiredVar)
	}
}

func TestRequiredMissing(t *testing.T) {
	var s Specification
	os.Clearenv()

	err := Process("env_config", &s)
	if err == nil {
		t.Error("no failure when missing required variable")
	}
}

func TestBlankDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "requiredvalue")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.DefaultVar)
	}

	if *s.SomePointerWithDefault != "foo2baz" {
		t.Errorf("expected %s, got %s", "foo2baz", *s.SomePointerWithDefault)
	}
}

func TestNonBlankDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEFAULTVAR", "nondefaultval")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "requiredvalue")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "nondefaultval" {
		t.Errorf("expected %s, got %s", "nondefaultval", s.DefaultVar)
	}
}

func TestExplicitBlankDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEFAULTVAR", "")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "")

	if err := Process("env_config", &s); err == nil {
		t.Error("no failure when missing required variable")
	}

	if s.DefaultVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.DefaultVar)
	}
}

func TestAlternateNameDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("BROKER", "betterbroker")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// This is also an inversion of the original test, since we no longer fallback
	// on the non-prefixed Alt version if the specified tag is not found.
	if s.NoPrefixDefault != "127.0.0.1" {
		t.Errorf("expected %q, got %q", "127.0.0.1", s.NoPrefixDefault)
	}

	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.NoPrefixDefault != "127.0.0.1" {
		t.Errorf("expected %q, got %q", "127.0.0.1", s.NoPrefixDefault)
	}
}

func TestRequiredDefault(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.RequiredDefault != "foo2bar" {
		t.Errorf("expected %q, got %q", "foo2bar", s.RequiredDefault)
	}
}

func TestPointerFieldBlank(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.SomePointer != nil {
		t.Errorf("expected <nil>, got %q", *s.SomePointer)
	}
}

func TestEmptyMapFieldOverride(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	os.Setenv("ENV_CONFIG_MAPFIELD", "")
	os.Setenv("ENV_CONFIG_EMPTY_MAPFIELD", "")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.MapField == nil {
		t.Error("expected map, got <nil>")
	}

	expMap := map[string]string{
		"one":   "two",
		"three": "four",
	}
	if !reflect.DeepEqual(s.MapField, expMap) {
		t.Errorf("expected map %+v, got map %+v", expMap, s.MapField)
	}

	if s.EmptyMapField != nil {
		t.Errorf("expected nil map, but got %+v", s.EmptyMapField)
	}
}

func TestMustProcess(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	MustProcess("env_config", &s)

	defer func() {
		if err := recover(); err != nil {
			return
		}

		t.Error("expected panic")
	}()
	m := make(map[string]string)
	MustProcess("env_config", &m)
}

func TestEmbeddedStruct(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "required")
	os.Setenv("ENV_CONFIG_ENABLED", "true")
	os.Setenv("ENV_CONFIG_EMBEDDEDPORT", "1234")
	os.Setenv("ENV_CONFIG_MULTIWORDVAR", "foo")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("ENV_CONFIG_MULTI_WITH_DIFFERENT_ALT", "baz")
	os.Setenv("ENV_CONFIG_EMBEDDED_WITH_ALT", "foobar")
	os.Setenv("ENV_CONFIG_SOMEPOINTER", "foobaz")
	os.Setenv("ENV_CONFIG_EMBEDDED_IGNORED", "was-not-ignored")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if !s.Enabled {
		t.Errorf("expected %v, got %v", true, s.Enabled)
	}
	if s.EmbeddedPort != 1234 {
		t.Errorf("expected %d, got %v", 1234, s.EmbeddedPort)
	}
	if s.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.MultiWordVar)
	}
	if s.Embedded.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.Embedded.MultiWordVar)
	}
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %s, got %s", "bar", s.MultiWordVarWithAlt)
	}
	if s.Embedded.MultiWordVarWithAlt != "baz" {
		t.Errorf("expected %s, got %s", "baz", s.Embedded.MultiWordVarWithAlt)
	}
	if s.EmbeddedAlt != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.EmbeddedAlt)
	}
	if *s.SomePointer != "foobaz" {
		t.Errorf("expected %s, got %s", "foobaz", *s.SomePointer)
	}
	if s.EmbeddedIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestDayDuration(t *testing.T) {
	var s struct {
		Days0  time.Duration `envconfig:"DAYS_0" default:"0d"`
		Days1  time.Duration `envconfig:"DAYS_1" default:"1d"`
		Days10 time.Duration `envconfig:"DAYS_10" default:"10d"`
	}

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.Days0 != 0 {
		t.Errorf("expected %d, got %s", 0, s.Days0)
	}

	if s.Days10 != 10*24*time.Hour {
		t.Errorf("expected %s, got %s", 10*24*time.Hour, s.Days1)
	}

	if s.Days1 != 1*24*time.Hour {
		t.Errorf("expected %s, got %s", 10*24*time.Hour, s.Days10)
	}
}

func TestInvalidDayDuration(t *testing.T) {

	badDays := []string{
		"1dd",
		"d",
		" d",
	}

	for _, badDay := range badDays {
		var s Specification
		os.Clearenv()
		os.Setenv("ENV_CONFIG_TIMEOUT", badDay)
		err := Process("env_config", &s)

		if err == nil {
			t.Errorf("expected an err!")
		}
	}
}

func TestEmbeddedButIgnoredStruct(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "required")
	os.Setenv("ENV_CONFIG_FIRSTEMBEDDEDBUTIGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_SECONDEMBEDDEDBUTIGNORED", "was-not-ignored")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if s.FirstEmbeddedButIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
	if s.SecondEmbeddedButIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestNonPointerFailsProperly(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "snap")

	err := Process("env_config", s)
	if err != ErrInvalidSpecification {
		t.Errorf("non-pointer should fail with ErrInvalidSpecification, was instead %s", err)
	}
}

func TestCustomValueFields(t *testing.T) {
	var s struct {
		Foo    string       `envconfig:"FOO"`
		Bar    bracketed    `envconfig:"BAR"`
		Baz    quoted       `envconfig:"BAZ"`
		Struct setterStruct `envconfig:"STRUCT"`
	}

	// Set would panic when the receiver is nil,
	// so make sure it has an initial value to replace.
	s.Baz = quoted{new(bracketed)}

	os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")
	os.Setenv("ENV_CONFIG_BAZ", "baz")
	os.Setenv("ENV_CONFIG_STRUCT", "inner")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if want := "foo"; s.Foo != want {
		t.Errorf("foo: got %#q, want %#q", s.Foo, want)
	}

	if want := "[bar]"; s.Bar.String() != want {
		t.Errorf("bar: got %#q, want %#q", s.Bar, want)
	}

	if want := `["baz"]`; s.Baz.String() != want {
		t.Errorf(`baz: got %#q, want %#q`, s.Baz, want)
	}

	if want := `setterstruct{"inner"}`; s.Struct.Inner != want {
		t.Errorf(`Struct.Inner: got %#q, want %#q`, s.Struct.Inner, want)
	}
}

func TestCustomPointerFields(t *testing.T) {
	var s struct {
		Foo    string        `envconfig:"FOO"`
		Bar    *bracketed    `envconfig:"BAR"`
		Baz    *quoted       `envconfig:"BAZ"`
		Struct *setterStruct `envconfig:"STRUCT"`
	}

	// Set would panic when the receiver is nil,
	// so make sure they have initial values to replace.
	s.Bar = new(bracketed)
	s.Baz = &quoted{new(bracketed)}

	os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")
	os.Setenv("ENV_CONFIG_BAZ", "baz")
	os.Setenv("ENV_CONFIG_STRUCT", "inner")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if want := "foo"; s.Foo != want {
		t.Errorf("foo: got %#q, want %#q", s.Foo, want)
	}

	if want := "[bar]"; s.Bar.String() != want {
		t.Errorf("bar: got %#q, want %#q", s.Bar, want)
	}

	if want := `["baz"]`; s.Baz.String() != want {
		t.Errorf(`baz: got %#q, want %#q`, s.Baz, want)
	}

	if want := `setterstruct{"inner"}`; s.Struct.Inner != want {
		t.Errorf(`Struct.Inner: got %#q, want %#q`, s.Struct.Inner, want)
	}
}

func TestEmptyPrefixUsesFieldNames(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("REQUIREDVAR", "foo")

	err := Process("", &s)
	if err != nil {
		t.Errorf("Process failed: %s", err)
	}

	if s.RequiredVar != "foo" {
		t.Errorf(
			`RequiredVar not populated correctly: expected "foo", got %q`,
			s.RequiredVar,
		)
	}
}

func TestNestedStructVarName(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "required")
	// The behaviour of this test has changed, as we explicitly expect the prefix
	// to be used consistently, we no longer check for the INNER without the prefix.
	val := "found only with prefixed name"
	os.Setenv("ENV_CONFIG_OUTER_INNER", val)
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if s.NestedSpecification.Property != val {
		t.Errorf("expected %s, got %s", val, s.NestedSpecification.Property)
	}
}

func TestTextUnmarshalerError(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	os.Setenv("ENV_CONFIG_DATETIME", "I'M NOT A DATE")

	err := Process("env_config", &s)

	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Datetime" {
		t.Errorf("expected %s, got %v", "Datetime", v.FieldName)
	}

	expectedLowLevelError := time.ParseError{
		Layout:     time.RFC3339,
		Value:      "I'M NOT A DATE",
		LayoutElem: "2006",
		ValueElem:  "I'M NOT A DATE",
	}

	if v.Err.Error() != expectedLowLevelError.Error() {
		t.Errorf("expected %s, got %s", expectedLowLevelError, v.Err)
	}
}

func TestBinaryUnmarshalerError(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	os.Setenv("ENV_CONFIG_URLPOINTER", "http://%41:8080/")

	err := Process("env_config", &s)

	v, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError, got %T %v", err, err)
	}
	if v.FieldName != "UrlPointer" {
		t.Errorf("expected %s, got %v", "UrlPointer", v.FieldName)
	}

	// To be compatible with go 1.5 and lower we should do a very basic check,
	// because underlying error message varies in go 1.5 and go 1.6+.

	ue, ok := v.Err.(*url.Error)
	if !ok {
		t.Errorf("expected error type to be \"*url.Error\", got %T", v.Err)
	}

	if ue.Op != "parse" {
		t.Errorf("expected error op to be \"parse\", got %q", ue.Op)
	}
}

func TestCheckDisallowedOnlyAllowed(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("UNRELATED_ENV_VAR", "true")
	err := CheckDisallowed("env_config", &s)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestCheckDisallowedMispelled(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_ZEBUG", "false")
	err := CheckDisallowed("env_config", &s)
	if experr := "unknown environment variable ENV_CONFIG_ZEBUG"; err.Error() != experr {
		t.Errorf("expected %s, got %s", experr, err)
	}
}

func TestCheckDisallowedIgnored(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_IGNORED", "false")
	err := CheckDisallowed("env_config", &s)
	if experr := "unknown environment variable ENV_CONFIG_IGNORED"; err.Error() != experr {
		t.Errorf("expected %s, got %s", experr, err)
	}
}

func TestErrorMessageForRequiredAltVar(t *testing.T) {
	var s struct {
		Foo string `envconfig:"BAR" required:"true"`
	}

	os.Clearenv()
	err := Process("env_config", &s)

	if err == nil {
		t.Error("no failure when missing required variable")
	}

	if !strings.Contains(err.Error(), " ENV_CONFIG_BAR ") {
		t.Errorf("expected error message to contain ENV_CONFIG_BAR, got \"%v\"", err)
	}
}

func TestErrorMessageForRequiredAltVarNoPrefix(t *testing.T) {
	var s struct {
		Foo string `envconfig:"BAR" required:"true"`
	}

	os.Clearenv()
	err := Process("", &s)

	if err == nil {
		t.Error("no failure when missing required variable")
	}

	if !strings.Contains(err.Error(), " BAR ") {
		t.Errorf("expected error message to contain BAR, got \"%v\"", err)
	}
}

func TestErrorMessageForRequiredAltVarNestedStruct(t *testing.T) {
	var s struct {
		Foo struct {
			Bar string `envconfig:"BAR" required:"true"`
		} `envconfig:"FOO" required:"true"`
	}

	os.Clearenv()
	err := Process("ENV_CONFIG", &s)

	if err == nil {
		t.Error("no failure when missing required variable")
	}

	if !strings.Contains(err.Error(), " ENV_CONFIG_FOO_BAR ") {
		t.Errorf("expected error message to contain ENV_CONFIG_FOO_BAR, got \"%v\"", err)
	}
}

func TestNonTaggedFields(t *testing.T) {
	var s struct {
		Foo string `envconfig:"FOO"`
		Bar string `default:"can you hear me?"`
		Baz string `required:"true"`
	}

	os.Clearenv()
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	os.Setenv("BAZ", "baz")

	err := Process("", &s)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if s.Foo != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.Foo)
	}
	if s.Bar != "" {
		t.Errorf("expected %s, got %s", "", s.Bar)
	}
	if s.Baz != "" {
		t.Errorf("expected %s, got %s", "", s.Baz)
	}
}

func TestNestedStructs(t *testing.T) {
	var s struct {
		Anonymous struct {
			Foo string `envconfig:"FOO"`
			Bar string
		}
		Named struct {
			Foz string `envconfig:"FOZ"`
			Baz string
		} `envconfig:"NAMED"`
		NamedButSkipped struct {
			Hello string `envconfig:"HELLO"`
		} `envconfig:"SKIPPED"`
	}

	os.Clearenv()
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	os.Setenv("NAMED_FOZ", "foz")
	os.Setenv("NAMED_BAZ", "baz")
	os.Setenv("HELLO", "hello")

	err := Process("", &s)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if s.Anonymous.Foo != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.Anonymous.Foo)
	}
	if s.Anonymous.Bar != "" {
		t.Errorf("expected %s, got %s", "", s.Anonymous.Bar)
	}
	if s.Named.Foz != "foz" {
		t.Errorf("expected %s, got %s", "", s.Named.Foz)
	}
	if s.Named.Baz != "" {
		t.Errorf("expected %s, got %s", "", s.Named.Baz)
	}
	if s.NamedButSkipped.Hello != "" {
		t.Errorf("expected %s, got %s", "", s.NamedButSkipped.Hello)
	}
}

type bracketed string

func (b *bracketed) Set(value string) error {
	*b = bracketed("[" + value + "]")
	return nil
}

func (b bracketed) String() string {
	return string(b)
}

// quoted is used to test the precedence of Decode over Set.
// The sole field is a flag.Value rather than a setter to validate that
// all flag.Value implementations are also Setter implementations.
type quoted struct{ flag.Value }

func (d quoted) Decode(value string) error {
	return d.Set(`"` + value + `"`)
}

type setterStruct struct {
	Inner string
}

func (ss *setterStruct) Set(value string) error {
	ss.Inner = fmt.Sprintf("setterstruct{%q}", value)
	return nil
}

func BenchmarkGatherInfo(b *testing.B) {
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("ENV_CONFIG_TIMEOUT", "2m")
	os.Setenv("ENV_CONFIG_ADMINUSERS", "John,Adam,Will")
	os.Setenv("ENV_CONFIG_MAGICNUMBERS", "5,10,20")
	os.Setenv("ENV_CONFIG_COLORCODES", "red:1,green:2,blue:3")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_TTL", "30")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	os.Setenv("ENV_CONFIG_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_OUTER_INNER", "iamnested")
	os.Setenv("ENV_CONFIG_AFTERNESTED", "after")
	os.Setenv("ENV_CONFIG_HONOR", "honor")
	os.Setenv("ENV_CONFIG_DATETIME", "2016-08-16T18:57:05Z")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_AUTO_SPLIT", "24")
	for i := 0; i < b.N; i++ {
		var s Specification
		gatherInfo("env_config", &s)
	}
}
