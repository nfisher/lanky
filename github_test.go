package main

import (
	"encoding/json"
	"strings"
	"testing"
)

const validPushRequst = `{
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

const validPingRequest = `{
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
	r := strings.NewReader(validPushRequst)
	gp := &GithubPushPayload{}

	dec := json.NewDecoder(r)
	err := dec.Decode(gp)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

func Test_should_process_valid_ping_correctly(t *testing.T) {
	r := strings.NewReader(validPingRequest)
	gp := &GithubPingPayload{}

	dec := json.NewDecoder(r)
	err := dec.Decode(gp)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	var expectedUrl Url = "https://api.github.com/repos/hailocab/releases-web/hooks/4314541"
	if gp.Hook.Url != expectedUrl {
		t.Fatal("gp.Hook.Url = %v, want %v", gp.Hook.Url, expectedUrl)
	}
}
