package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Url string

type User struct {
	Name     string
	Email    string
	Username string
}

type Organization struct {
	Login            string
	Id               int
	Url              Url
	ReposUrl         Url `json:"repos_url"`
	EventsUrl        Url `json:"events_url"`
	MembersUrl       Url `json:"members_url"`
	PublicMembersUrl Url `json:"publicMembers_url"`
	AvatarUrl        Url `json:"avatar_url"`
	Description      string
}

type Sender struct {
	Login             string
	Id                int
	AvatarUrl         Url `json:"avatar_url"`
	GravatarId        string
	Url               Url
	HtmlUrl           Url `json:"html_url"`
	FollowersUrl      Url `json:"followers_url"`
	FollowingUrl      Url `json:"following_url"`
	GistsUrl          Url `json:"gists_url"`
	StarredUrl        Url `json:"starred_url"`
	SubscriptionsUrl  Url `json:"subscriptions_url"`
	OrganizationsUrl  Url `json:"organizations_url"`
	ReposUrl          Url `json:"repos_url"`
	EventsUrl         Url `json:"events_url"`
	ReceivedEventsUrl Url `json:"receivedEvents_url"`
	Type              string
	SiteAdmin         bool
}

type Commit struct {
	Id        string
	Distinct  bool
	Message   string
	Timestamp time.Time
	Url       Url
	Author    User
	Committer User
	Added     []string
	Removed   []string
	Modified  []string
}

type Repository struct {
	Id       int
	Name     string
	FullName string `json:"full_name"`
	Owner    struct {
		Name  string
		Email string
	}
	Private          bool
	HtmlUrl          Url `json:"html_url"`
	Description      string
	Fork             bool
	Url              Url
	ForksUrl         Url `json:"forks_url"`
	KeysUrl          Url `json:"keys_url"`
	CollaboratorsUrl Url `json:"collaborators_url"`
	IssueEventsUrl   Url `json:"issue_events_url"`
	EventsUrl        Url `json:"events_url"`
	AssigneesUrl     Url `json:"assignees_url"`
	BranchesUrl      Url `json:"branches_url"`
	TagsUrl          Url `json:"tags_url"`
	BlobsUrl         Url `json:"blobs_url"`
	GitTagsUrl       Url `json:"git_tags_url"`
	GitRefsUrl       Url `json:"git_refs_url"`
	TreesUrl         Url `json:"trees_url"`
	StatusesUrl      Url `json:"statuses_url"`
	LanguagesUrl     Url `json:"languages_url"`
	StargazersUrl    Url `json:"stargazers_url"`
	ContributorsUrl  Url `json:"contributors_url"`
	SubscribersUrl   Url `json:"subscribers_url"`
	SubscriptionUrl  Url `json:"subscription_url"`
	CommitsUrl       Url `json:"commits_url"`
	GitCommitsUrl    Url `json:"git_commits_url"`
	CommentsUrl      Url `json:"comments_url"`
	IssueCommentUrl  Url `json:"issue_comment_url"`
	ContentsUrl      Url `json:"contents_url"`
	CompareUrl       Url `json:"compare_url"`
	MergesUrl        Url `json:"merges_url"`
	ArchiveUrl       Url `json:"archive_url"`
	DownloadsUrl     Url `json:"downloads_url"`
	IssuesUrl        Url `json:"issues_url"`
	PullsUrl         Url `json:"pulls_url"`
	MilestonesUrl    Url `json:"milestones_url"`
	NotificationsUrl Url `json:"notifications_url"`
	LabelsUrl        Url `json:"labels_url"`
	ReleasesUrl      Url `json:"releases_url"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	PushedAt         time.Time
	GitUrl           Url `json:"git_url"`
	SshUrl           Url `json:"ssh_url"`
	CloneUrl         Url `json:"clone_url"`
	SvnUrl           Url `json:"svn_url"`
	Homepage         Url `json:"homepage"`
	Size             int
	StargazersCount  int
	WatchersCount    int
	Language         string
	HasIssues        bool
	HasDownloads     bool
	HasWiki          bool
	HasPages         bool
	ForksCount       int
	MirrorsUrl       Url `json:"mirrors_url"`
	OpenIssuesCount  int
	Forks            int
	OpenIssues       int
	Watchers         int
	DefaultBranch    string
	Stargazers       int
	MasterBranch     string
	Organization     string
}

type Repositories []Repository

func (r Repositories) Len() int      { return len(r) }
func (r Repositories) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type ByFullName struct{ Repositories }

func (r ByFullName) Less(i, j int) bool {
	return r.Repositories[i].FullName < r.Repositories[j].FullName
}

type GithubPushPayload struct {
	Ref          string
	Before       string
	After        string
	Created      bool
	Deleted      bool
	Forced       bool
	BaseRef      *string
	Compare      Url
	Commits      []Commit
	HeadCommit   Commit
	Repository   Repository
	Pusher       User
	Organization Organization
}

type GithubPingPayload struct {
	Zen    string
	HookId int
	Hook   struct {
		Url     Url
		TestUrl Url `json:"test_url"`
		PingUrl Url `json:"ping_url"`
		Id      int
		Name    string
		Active  bool
		Events  []string
		Config  struct {
			Url         Url
			Secret      string
			ContentType string
		}
		LastResponse struct {
			Code    int
			Status  string
			Message string
		}
		UpdatedAt time.Time
		CreatedAt time.Time
	}
	Repository Repository
	Sender     Sender
}

func NewGithub(config *Config) (client *GithubClient) {
	if config == nil || config.Github == nil || config.Github.Token == "" {
		return nil
	}

	oa2conf := &oauth2.Config{
		Scopes:   []string{},
		Endpoint: github.Endpoint,
	}

	token := &oauth2.Token{
		AccessToken: config.Github.Token,
	}

	wc := oa2conf.Client(oauth2.NoContext, token)

	return &GithubClient{
		config,
		wc,
	}
}

type GithubClient struct {
	*Config
	WebClient
}

var linkRegex *regexp.Regexp = regexp.MustCompile(`<([^>]*)>;\s+rel="([^"]+)"`)

func (gc *GithubClient) Pagination(link, rel string) string {
	s := linkRegex.FindAllStringSubmatch(link, -1)
	for i := range s {
		if len(s[i]) == 3 {
			if s[i][2] == rel {
				return s[i][1]
			}
		}
	}
	return ""
}

const (
	relNext = "next"
	relLast = "last"

	linkHeader = "Link"
)

func (gc *GithubClient) Repositories(org string, repos *Repositories) (err error) {
	c := cap(*repos)
	repoPath := fmt.Sprintf("https://api.github.com/orgs/%v/repos?per_page=%v", org, c)
	// TODO: (NF 2015-04-05) put a reasonable limit of say 2000 repositories
	for {
		repositories := make(Repositories, 0, c)
		resp, err := gc.WebClient.Get(repoPath)
		if err != nil {
			return err
		}

		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&repositories)
		resp.Body.Close()
		if err != nil {
			return err
		}

		*repos = append(*repos, repositories...)

		link := resp.Header.Get(linkHeader)
		repoPath = gc.Pagination(link, relNext)
		if repoPath == "" {
			break
		}
	}

	sort.Sort(ByFullName{*repos})

	return nil
}
