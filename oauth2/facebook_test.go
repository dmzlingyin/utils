package oauth2

import (
	"context"
	"golang.org/x/oauth2"
	"io"
	"testing"
)

func TestFacebookToken(t *testing.T) {
	f := NewFacebook()
	token := &oauth2.Token{
		AccessToken: "EAAIBCJ7M9rgBAKOdClNoC0UaLtDuXP5IJhZCmj6cEmXH0DWUGqMcHM8TZAAU8YDoir6bolDHjwBdNJmlshzER1brAmqXRZAZCXlAJfwR75I87u2PABfgB72aV1SmkclKk2s1yOGBy4rrMsAkZBNZARleis7RqBjiMP9uuCba673xfxqPLurSYrAoj8GClH3ZBC2lkFu4Jj9JnZBMT8bWqskOEDVM2GbZAtusEbROF6iNxHFJ9T1Y4SOBP",
		TokenType:   "Bearer",
	}
	res, err := f.cfg.Client(context.TODO(), token).Get(FacebookUserURL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(string(body))
}

func TestFacebookCode(t *testing.T) {
	code := "AQCL0WhlRDwUR3Bn7a6ze_a-CntdvrtCN6lUzH3B5S3ZXRAPJltz-7qU85xrmPqgtHUFeTd7MUa7_JPMvrIXY3SpS9vMa1V5uYZXNpmKqjFjQ3wBXlkHswM6vAlUDfrC1y6tsn4oGCCcoJc3UkNApNiK6Ixt5ZH9xTxW8mO0uqLEVme1y-BNiJ805F_RMLBjpB0aXYRItc453MTQVbLMeYVKTy4NdcMn60LsW6DvtZvIHs5J0nm5vHzkWhfTR5ngDbuqzogT1fCQ6K0aps6ur6KnNs-uBfpsMClPtnhLj0fxIgdRR-pBnLg_L3j5ij873lflV3rqkNk_efxKLtOkQO-boj512rwsKvOj93E0QvzCRn7kv1HR6gcepIuXoKRihYY"
	f := NewFacebook()
	token, user, err := f.Authorize(context.Background(), code)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %+v, user: %+v", token, user)
}
