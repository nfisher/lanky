package main

import "time"

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
	ReposUrl         Url
	EventsUrl        Url
	MembersUrl       Url
	PublicMembersUrl Url
	AvatarUrl        Url
	Description      string
}

type Sender struct {
	Login             string
	Id                int
	AvatarUrl         Url
	GravatarId        string
	Url               Url
	HtmlUrl           Url
	FollowersUrl      Url
	FollowingUrl      Url
	GistsUrl          Url
	StarredUrl        Url
	SubscriptionsUrl  Url
	OrganizationsUrl  Url
	ReposUrl          Url
	EventsUrl         Url
	ReceivedEventsUrl Url
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
	FullName string
	Owner    struct {
		Name  string
		Email string
	}
	Private          bool
	HtmlUrl          Url
	Description      string
	Fork             bool
	Url              Url
	ForksUrl         Url
	KeysUrl          Url
	CollaboratorsUrl Url
	IssueEventsUrl   Url
	EventsUrl        Url
	AssigneesUrl     Url
	BranchesUrl      Url
	TagsUrl          Url
	BlobsUrl         Url
	GitTagsUrl       Url
	GitRefsUrl       Url
	TreesUrl         Url
	StatusesUrl      Url
	LanguagesUrl     Url
	StargazersUrl    Url
	ContributorsUrl  Url
	SubscribersUrl   Url
	SubscriptionUrl  Url
	CommitsUrl       Url
	GitCommitsUrl    Url
	CommentsUrl      Url
	IssueCommentUrl  Url
	ContentsUrl      Url
	CompareUrl       Url
	MergesUrl        Url
	ArchiveUrl       Url
	DownloadsUrl     Url
	IssuesUrl        Url
	PullsUrl         Url
	MilestonesUrl    Url
	NotificationsUrl Url
	LabelsUrl        Url
	ReleasesUrl      Url
	CreatedAt        time.Time
	UpdatedAt        time.Time
	PushedAt         time.Time
	GitUrl           Url
	SshUrl           Url
	CloneUrl         Url
	SvnUrl           Url
	Homepage         Url
	Size             int
	StargazersCount  int
	WatchersCount    int
	Language         string
	HasIssues        bool
	HasDownloads     bool
	HasWiki          bool
	HasPages         bool
	ForksCount       int
	MirrorsUrl       Url
	OpenIssuesCount  int
	Forks            int
	OpenIssues       int
	Watchers         int
	DefaultBranch    string
	Stargazers       int
	MasterBranch     string
	Organization     string
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
		TestUrl Url
		PingUrl Url
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
