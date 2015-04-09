package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

const validHookResponse = `[
  {
    "id": 1,
    "url": "https://api.github.com/repos/octocat/Hello-World/hooks/1",
    "test_url": "https://api.github.com/repos/octocat/Hello-World/hooks/1/test",
    "ping_url": "https://api.github.com/repos/octocat/Hello-World/hooks/1/pings",
    "name": "web",
    "events": [
      "push",
      "pull_request"
    ],
    "active": true,
    "config": {
      "url": "http://example.com/webhook",
      "content_type": "json"
    },
    "updated_at": "2011-09-06T20:39:23Z",
    "created_at": "2011-09-06T17:26:27Z"
  }
]`

const validRepositoriesResponse = `[
  {
    "id": 1296269,
    "owner": {
      "login": "octocat",
      "id": 1,
      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
      "gravatar_id": "",
      "url": "https://api.github.com/users/octocat",
      "html_url": "https://github.com/octocat",
      "followers_url": "https://api.github.com/users/octocat/followers",
      "following_url": "https://api.github.com/users/octocat/following{/other_user}",
      "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
      "organizations_url": "https://api.github.com/users/octocat/orgs",
      "repos_url": "https://api.github.com/users/octocat/repos",
      "events_url": "https://api.github.com/users/octocat/events{/privacy}",
      "received_events_url": "https://api.github.com/users/octocat/received_events",
      "type": "User",
      "site_admin": false
    },
    "name": "Hello-World",
    "full_name": "octocat/Hello-World",
    "description": "This your first repo!",
    "private": false,
    "fork": false,
    "url": "https://api.github.com/repos/octocat/Hello-World",
    "html_url": "https://github.com/octocat/Hello-World",
    "clone_url": "https://github.com/octocat/Hello-World.git",
    "git_url": "git://github.com/octocat/Hello-World.git",
    "ssh_url": "git@github.com:octocat/Hello-World.git",
    "svn_url": "https://svn.github.com/octocat/Hello-World",
    "mirror_url": "git://git.example.com/octocat/Hello-World",
    "homepage": "https://github.com",
    "language": null,
    "forks_count": 9,
    "stargazers_count": 80,
    "watchers_count": 80,
    "size": 108,
    "default_branch": "master",
    "open_issues_count": 0,
    "has_issues": true,
    "has_wiki": true,
    "has_pages": false,
    "has_downloads": true,
    "pushed_at": "2011-01-26T19:06:43Z",
    "created_at": "2011-01-26T19:01:12Z",
    "updated_at": "2011-01-26T19:14:43Z",
    "permissions": {
      "admin": false,
      "push": false,
      "pull": true
    }
  }
]`

const validPushResponse = `{
  "ref": "refs/heads/master",
  "before": "98631d4912c3e4dbad586ea01a00274d364e0745",
  "after": "ebe220cce16e1d9ff50b7bf0de5033ff89c4ed81",
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/hailocab/releases-web/compare/98631d4912c3...ebe220cce16e",
  "commits": [
    {
      "id": "ebe220cce16e1d9ff50b7bf0de5033ff89c4ed81",
      "distinct": true,
      "message": "[#INF-3156] Add listing of provisioned services and stub for modification of existing services. :@nfisher",
      "timestamp": "2015-03-30T13:55:47+01:00",
      "url": "https://github.com/hailocab/releases-web/commit/ebe220cce16e1d9ff50b7bf0de5033ff89c4ed81",
      "author": {
        "name": "Nathan Fisher",
        "email": "nfisher@junctionbox.ca",
        "username": "nfisher"
      },
      "committer": {
        "name": "Nathan Fisher",
        "email": "nfisher@junctionbox.ca",
        "username": "nfisher"
      },
      "added": [
        "media/hailo-logo-black.svg",
        "releases/modify.html",
        "src/release_list_view.js",
        "src/release_modify_view.js"
      ],
      "removed": [

      ],
      "modified": [
        "SpecRunner.html",
        "_layouts/default.html",
        "css/main.scss",
        "login.html",
        "releases/index.html",
        "releases/service.html",
        "src/main.js",
        "src/routes.js",
        "src/traffic_list_view.js",
        "traffic/index.html",
        "traffic/modify.html"
      ]
    }
  ],
  "head_commit": {
    "id": "ebe220cce16e1d9ff50b7bf0de5033ff89c4ed81",
    "distinct": true,
    "message": "[#INF-3156] Add listing of provisioned services and stub for modification of existing services. :@nfisher",
    "timestamp": "2015-03-30T13:55:47+01:00",
    "url": "https://github.com/hailocab/releases-web/commit/ebe220cce16e1d9ff50b7bf0de5033ff89c4ed81",
    "author": {
      "name": "Nathan Fisher",
      "email": "nfisher@junctionbox.ca",
      "username": "nfisher"
    },
    "committer": {
      "name": "Nathan Fisher",
      "email": "nfisher@junctionbox.ca",
      "username": "nfisher"
    },
    "added": [
      "media/hailo-logo-black.svg",
      "releases/modify.html",
      "src/release_list_view.js",
      "src/release_modify_view.js"
    ],
    "removed": [

    ],
    "modified": [
      "SpecRunner.html",
      "_layouts/default.html",
      "css/main.scss",
      "login.html",
      "releases/index.html",
      "releases/service.html",
      "src/main.js",
      "src/routes.js",
      "src/traffic_list_view.js",
      "traffic/index.html",
      "traffic/modify.html"
    ]
  },
  "repository": {
    "id": 28084179,
    "name": "releases-web",
    "full_name": "hailocab/releases-web",
    "owner": {
      "name": "hailocab",
      "email": ""
    },
    "private": true,
    "html_url": "https://github.com/hailocab/releases-web",
    "description": "Unleash the hounds... err services.",
    "fork": false,
    "url": "https://github.com/hailocab/releases-web",
    "forks_url": "https://api.github.com/repos/hailocab/releases-web/forks",
    "keys_url": "https://api.github.com/repos/hailocab/releases-web/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/hailocab/releases-web/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/hailocab/releases-web/teams",
    "hooks_url": "https://api.github.com/repos/hailocab/releases-web/hooks",
    "issue_events_url": "https://api.github.com/repos/hailocab/releases-web/issues/events{/number}",
    "events_url": "https://api.github.com/repos/hailocab/releases-web/events",
    "assignees_url": "https://api.github.com/repos/hailocab/releases-web/assignees{/user}",
    "branches_url": "https://api.github.com/repos/hailocab/releases-web/branches{/branch}",
    "tags_url": "https://api.github.com/repos/hailocab/releases-web/tags",
    "blobs_url": "https://api.github.com/repos/hailocab/releases-web/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/hailocab/releases-web/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/hailocab/releases-web/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/hailocab/releases-web/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/hailocab/releases-web/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/hailocab/releases-web/languages",
    "stargazers_url": "https://api.github.com/repos/hailocab/releases-web/stargazers",
    "contributors_url": "https://api.github.com/repos/hailocab/releases-web/contributors",
    "subscribers_url": "https://api.github.com/repos/hailocab/releases-web/subscribers",
    "subscription_url": "https://api.github.com/repos/hailocab/releases-web/subscription",
    "commits_url": "https://api.github.com/repos/hailocab/releases-web/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/hailocab/releases-web/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/hailocab/releases-web/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/hailocab/releases-web/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/hailocab/releases-web/contents/{+path}",
    "compare_url": "https://api.github.com/repos/hailocab/releases-web/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/hailocab/releases-web/merges",
    "archive_url": "https://api.github.com/repos/hailocab/releases-web/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/hailocab/releases-web/downloads",
    "issues_url": "https://api.github.com/repos/hailocab/releases-web/issues{/number}",
    "pulls_url": "https://api.github.com/repos/hailocab/releases-web/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/hailocab/releases-web/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/hailocab/releases-web/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/hailocab/releases-web/labels{/name}",
    "releases_url": "https://api.github.com/repos/hailocab/releases-web/releases{/id}",
    "created_at": 1418729544,
    "updated_at": "2015-03-30T13:00:14Z",
    "pushed_at": 1427720414,
    "git_url": "git://github.com/hailocab/releases-web.git",
    "ssh_url": "git@github.com:hailocab/releases-web.git",
    "clone_url": "https://github.com/hailocab/releases-web.git",
    "svn_url": "https://github.com/hailocab/releases-web",
    "homepage": null,
    "size": 440,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": "JavaScript",
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 0,
    "forks": 0,
    "open_issues": 0,
    "watchers": 0,
    "default_branch": "master",
    "stargazers": 0,
    "master_branch": "master",
    "organization": "hailocab"
  },
  "pusher": {
    "name": "nfisher",
    "email": "nfisher@junctionbox.ca"
  },
  "organization": {
    "login": "hailocab",
    "id": 560650,
    "url": "https://api.github.com/orgs/hailocab",
    "repos_url": "https://api.github.com/orgs/hailocab/repos",
    "events_url": "https://api.github.com/orgs/hailocab/events",
    "members_url": "https://api.github.com/orgs/hailocab/members{/member}",
    "public_members_url": "https://api.github.com/orgs/hailocab/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/560650?v=3",
    "description": ""
  },
  "sender": {
    "login": "nfisher",
    "id": 18616,
    "avatar_url": "https://avatars.githubusercontent.com/u/18616?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/nfisher",
    "html_url": "https://github.com/nfisher",
    "followers_url": "https://api.github.com/users/nfisher/followers",
    "following_url": "https://api.github.com/users/nfisher/following{/other_user}",
    "gists_url": "https://api.github.com/users/nfisher/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/nfisher/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/nfisher/subscriptions",
    "organizations_url": "https://api.github.com/users/nfisher/orgs",
    "repos_url": "https://api.github.com/users/nfisher/repos",
    "events_url": "https://api.github.com/users/nfisher/events{/privacy}",
    "received_events_url": "https://api.github.com/users/nfisher/received_events",
    "type": "User",
    "site_admin": false
  }
}`

const validPingResponse = `{
  "zen": "Half measures are as bad as nothing at all.",
  "hook_id": 4314541,
  "hook": {
    "url": "https://api.github.com/repos/hailocab/releases-web/hooks/4314541",
    "test_url": "https://api.github.com/repos/hailocab/releases-web/hooks/4314541/test",
    "ping_url": "https://api.github.com/repos/hailocab/releases-web/hooks/4314541/pings",
    "id": 4314541,
    "name": "web",
    "active": true,
    "events": [
      "push"
    ],
    "config": {
      "url": "http://janky.hailoweb.com:9393/_github",
      "secret": "********",
      "content_type": "json"
    },
    "last_response": {
      "code": null,
      "status": "unused",
      "message": null
    },
    "updated_at": "2015-03-11T17:48:16Z",
    "created_at": "2015-03-11T17:48:16Z"
  },
  "repository": {
    "id": 28084179,
    "name": "releases-web",
    "full_name": "hailocab/releases-web",
    "owner": {
      "login": "hailocab",
      "id": 560650,
      "avatar_url": "https://avatars.githubusercontent.com/u/560650?v=3",
      "gravatar_id": "",
      "url": "https://api.github.com/users/hailocab",
      "html_url": "https://github.com/hailocab",
      "followers_url": "https://api.github.com/users/hailocab/followers",
      "following_url": "https://api.github.com/users/hailocab/following{/other_user}",
      "gists_url": "https://api.github.com/users/hailocab/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/hailocab/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/hailocab/subscriptions",
      "organizations_url": "https://api.github.com/users/hailocab/orgs",
      "repos_url": "https://api.github.com/users/hailocab/repos",
      "events_url": "https://api.github.com/users/hailocab/events{/privacy}",
      "received_events_url": "https://api.github.com/users/hailocab/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "private": true,
    "html_url": "https://github.com/hailocab/releases-web",
    "description": "Unleash the hounds... err services.",
    "fork": false,
    "url": "https://api.github.com/repos/hailocab/releases-web",
    "forks_url": "https://api.github.com/repos/hailocab/releases-web/forks",
    "keys_url": "https://api.github.com/repos/hailocab/releases-web/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/hailocab/releases-web/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/hailocab/releases-web/teams",
    "hooks_url": "https://api.github.com/repos/hailocab/releases-web/hooks",
    "issue_events_url": "https://api.github.com/repos/hailocab/releases-web/issues/events{/number}",
    "events_url": "https://api.github.com/repos/hailocab/releases-web/events",
    "assignees_url": "https://api.github.com/repos/hailocab/releases-web/assignees{/user}",
    "branches_url": "https://api.github.com/repos/hailocab/releases-web/branches{/branch}",
    "tags_url": "https://api.github.com/repos/hailocab/releases-web/tags",
    "blobs_url": "https://api.github.com/repos/hailocab/releases-web/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/hailocab/releases-web/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/hailocab/releases-web/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/hailocab/releases-web/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/hailocab/releases-web/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/hailocab/releases-web/languages",
    "stargazers_url": "https://api.github.com/repos/hailocab/releases-web/stargazers",
    "contributors_url": "https://api.github.com/repos/hailocab/releases-web/contributors",
    "subscribers_url": "https://api.github.com/repos/hailocab/releases-web/subscribers",
    "subscription_url": "https://api.github.com/repos/hailocab/releases-web/subscription",
    "commits_url": "https://api.github.com/repos/hailocab/releases-web/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/hailocab/releases-web/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/hailocab/releases-web/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/hailocab/releases-web/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/hailocab/releases-web/contents/{+path}",
    "compare_url": "https://api.github.com/repos/hailocab/releases-web/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/hailocab/releases-web/merges",
    "archive_url": "https://api.github.com/repos/hailocab/releases-web/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/hailocab/releases-web/downloads",
    "issues_url": "https://api.github.com/repos/hailocab/releases-web/issues{/number}",
    "pulls_url": "https://api.github.com/repos/hailocab/releases-web/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/hailocab/releases-web/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/hailocab/releases-web/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/hailocab/releases-web/labels{/name}",
    "releases_url": "https://api.github.com/repos/hailocab/releases-web/releases{/id}",
    "created_at": "2014-12-16T11:32:24Z",
    "updated_at": "2015-03-09T15:20:30Z",
    "pushed_at": "2015-03-09T15:20:30Z",
    "git_url": "git://github.com/hailocab/releases-web.git",
    "ssh_url": "git@github.com:hailocab/releases-web.git",
    "clone_url": "https://github.com/hailocab/releases-web.git",
    "svn_url": "https://github.com/hailocab/releases-web",
    "homepage": null,
    "size": 328,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": "JavaScript",
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 0,
    "forks": 0,
    "open_issues": 0,
    "watchers": 0,
    "default_branch": "master"
  },
  "sender": {
    "login": "conanthedeployer",
    "id": 847044,
    "avatar_url": "https://avatars.githubusercontent.com/u/847044?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/conanthedeployer",
    "html_url": "https://github.com/conanthedeployer",
    "followers_url": "https://api.github.com/users/conanthedeployer/followers",
    "following_url": "https://api.github.com/users/conanthedeployer/following{/other_user}",
    "gists_url": "https://api.github.com/users/conanthedeployer/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/conanthedeployer/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/conanthedeployer/subscriptions",
    "organizations_url": "https://api.github.com/users/conanthedeployer/orgs",
    "repos_url": "https://api.github.com/users/conanthedeployer/repos",
    "events_url": "https://api.github.com/users/conanthedeployer/events{/privacy}",
    "received_events_url": "https://api.github.com/users/conanthedeployer/received_events",
    "type": "User",
    "site_admin": false
  }
}`

func Test_should_process_valid_push_correctly(t *testing.T) {
	r := strings.NewReader(validPushResponse)
	gp := &GithubPushPayload{}

	dec := json.NewDecoder(r)
	err := dec.Decode(gp)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

func Test_should_process_valid_ping_correctly(t *testing.T) {
	r := strings.NewReader(validPingResponse)
	gp := &GithubPingPayload{}

	dec := json.NewDecoder(r)
	err := dec.Decode(gp)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	var expectedUrl Url = "https://api.github.com/repos/hailocab/releases-web/hooks/4314541"
	if gp.Hook.Url != expectedUrl {
		t.Fatalf("gp.Hook.Url = %v, want %v", gp.Hook.Url, expectedUrl)
	}
}

const validNextPageHeader = `<https://api.github.com/organizations/560650/repos?per_page=100&page=2>; rel="next", <https://api.github.com/organizations/560650/repos?per_page=100&page=8>; rel="last"`

var pagination = []struct {
	rel      string
	expected string
}{
	{"next", "https://api.github.com/organizations/560650/repos?per_page=100&page=2"},
	{"last", "https://api.github.com/organizations/560650/repos?per_page=100&page=8"},
	{"first", ""},
}

func Test_Pagination(t *testing.T) {
	gc := &GithubClient{}

	for _, tt := range pagination {
		actual := gc.Pagination(validNextPageHeader, tt.rel)

		if actual != tt.expected {
			t.Fatalf("gc.Pagination(header,\"%v\") = %v, want %v", tt.rel, actual, tt.expected)
		}
	}
}

var githubConfigTable = []struct {
	config *Config
}{
	{nil},
	{&Config{}},
	{&Config{Github: &Github{Token: ""}}},
}

func Test_NewGithub_with_invalid_config(t *testing.T) {
	for _, tt := range githubConfigTable {
		gc := NewGithub(tt.config)
		if gc != nil {
			t.Fatalf("NewGithub(%v) = %v, want nil", tt.config, gc)
		}
	}
}

func Test_NewGithub_with_valid_token(t *testing.T) {
	config := &Config{Github: &Github{Token: "secret"}}
	gc := NewGithub(config)
	if gc == nil {
		t.Fatalf("NewGithub(%v) = nil, want &GithubClient{}", config)
	}
}

type TestClient struct {
	responses []string
	urls      []string
}

func newClient() *TestClient {
	return &TestClient{
		responses: make([]string, 0, 8),
		urls:      make([]string, 0, 8),
	}
}

type closer struct {
	*strings.Reader
}

func (c *closer) Close() error {
	return nil
}

func (tc *TestClient) Get(url string) (*http.Response, error) {
	if len(tc.responses) > 0 {
		cur := &closer{strings.NewReader(tc.responses[0])}
		tc.responses = append(tc.responses[:0], tc.responses[:1]...)
		resp := &http.Response{
			Body: cur,
		}

		tc.urls = append(tc.urls, url)

		return resp, nil
	}

	return nil, errors.New("No response specified")
}

func (tc *TestClient) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return nil, nil
}

func Test_ListHooks_with_connection_error_should_return_error(t *testing.T) {
	tc := newClient()

	gc := &GithubClient{
		WebClient: tc,
	}

	hooks := make(Hooks, 0, 10)
	err := gc.ListHooks("octocat/Hello-World", &hooks)
	if err == nil {
		t.Fatal("err = nil, want error")
	}
}

func Test_ListHooks_with_json_parsing_error_should_return_error(t *testing.T) {
	tc := newClient()
	tc.responses = append(tc.responses, validHookResponse[:len(validHookResponse)-2])
	gc := &GithubClient{
		WebClient: tc,
	}

	hooks := make(Hooks, 0, 10)
	err := gc.ListHooks("octocat/Hello-World", &hooks)
	if err == nil {
		t.Fatal("err = nil, want error")
	}

	expectedMsg := "unexpected EOF"
	if err.Error() != expectedMsg {
		t.Fatalf("err.Error() = %v, want %v", err.Error(), expectedMsg)
	}
}

func Test_ListHooks_with_valid_single_response(t *testing.T) {
	tc := newClient()
	tc.responses = append(tc.responses, validHookResponse)
	gc := &GithubClient{
		WebClient: tc,
	}

	hooks := make(Hooks, 0, 10)
	err := gc.ListHooks("octocat/Hello-World", &hooks)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	expectedLen := 1
	if len(tc.urls) != expectedLen {
		t.Fatalf("len(tc.urls) = %v, want %v", len(tc.urls), expectedLen)
	}

	expectedUrl := "https://api.github.com/repos/octocat/Hello-World/hooks?per_page=10"
	if tc.urls[0] != expectedUrl {
		t.Fatalf("tc.urls[0] = %v, want %v", tc.urls[0], expectedUrl)
	}

	expectedLen = 1
	if len(hooks) != expectedLen {
		t.Fatalf("len(hooks) = %v, want %v", len(hooks), expectedLen)
	}

	expecedContentType := "json"
	if hooks[0].Config.ContentType != expecedContentType {
		t.Fatalf("hooks[0].Config.ContentType = %v, want %v", hooks[0].Config.ContentType, expecedContentType)
	}
}

func Test_ListRepositories_with_connection_error_should_return_error(t *testing.T) {
	tc := newClient()
	gc := &GithubClient{
		WebClient: tc,
	}

	repos := make(Repositories, 0, 1)
	err := gc.ListRepositories("hailocab", &repos)
	if err == nil {
		t.Fatal("err = nil, want error")
	}
}

func Test_ListRepositories_with_invalid_json_response_should_return_error(t *testing.T) {
	tc := newClient()
	tc.responses = append(tc.responses, validRepositoriesResponse[:len(validRepositoriesResponse)-2])
	gc := &GithubClient{
		WebClient: tc,
	}

	repos := make(Repositories, 0, 1)
	err := gc.ListRepositories("hailocab", &repos)
	if err == nil {
		t.Fatal("err = nil, want error")
	}

	expectedMsg := "unexpected EOF"
	if err.Error() != expectedMsg {
		t.Fatalf("err.Error() = %v, want %v", err.Error(), expectedMsg)
	}
}

func Test_ListRepositories_with_single_valid_json_response_should_return_repository_list(t *testing.T) {
	tc := newClient()
	tc.responses = append(tc.responses, validRepositoriesResponse)
	gc := &GithubClient{
		WebClient: tc,
	}

	repos := make(Repositories, 0, 1)
	err := gc.ListRepositories("hailocab", &repos)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	expectedLen := 1
	if len(repos) != expectedLen {
		t.Fatalf("len(repos) = %v, want %v", len(repos), expectedLen)
	}

	expectedName := "octocat/Hello-World"
	if repos[0].FullName != expectedName {
		t.Fatalf(".FullName = %v, want %v", repos[0].FullName, expectedName)
	}
}
