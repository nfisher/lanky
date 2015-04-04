package main

import (
	"strings"
	"testing"
	"time"
)

const validTrayFeed = `<Projects><Project webUrl="http://ci.jenkins-ci.org/job/core_selenium-test/" name="core_selenium-test" lastBuildLabel="18" lastBuildTime="2012-10-27T01:55:14Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_backend-merge-all-repo/" name="infra_backend-merge-all-repo" lastBuildLabel="162" lastBuildTime="2015-03-29T03:23:00Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_backend-war-size-tracker/" name="infra_backend-war-size-tracker" lastBuildLabel="903" lastBuildTime="2015-04-01T04:35:58Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_commit_history_generation/" name="infra_commit_history_generation" lastBuildLabel="695" lastBuildTime="2015-04-01T18:48:00Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_extension-indexer/" name="infra_extension-indexer" lastBuildLabel="152" lastBuildTime="2015-03-28T15:48:00Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_github_repository_list/" name="infra_github_repository_list" lastBuildLabel="1511" lastBuildTime="2015-04-01T04:53:00Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_plugin_changes_report/" name="infra_plugin_changes_report" lastBuildLabel="359" lastBuildTime="2015-03-30T08:00:09Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_plugins_svn_to_git/" name="infra_plugins_svn_to_git" lastBuildLabel="768" lastBuildTime="2010-11-21T16:03:50Z" lastBuildStatus="Unknown" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_svnsync/" name="infra_svnsync" lastBuildLabel="21243" lastBuildTime="2011-02-06T18:31:36Z" lastBuildStatus="Unknown" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/infra_sync_maven-hpi-plugin_www/" name="infra_sync_maven-hpi-plugin_www" lastBuildLabel="507" lastBuildTime="2015-04-01T14:14:00Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/jenkins_pom/" name="jenkins_pom" lastBuildLabel="292" lastBuildTime="2015-03-29T20:27:00Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/jenkins_ui-changes_branch/" name="jenkins_ui-changes_branch" lastBuildLabel="33" lastBuildTime="2012-10-11T17:51:30Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/lib-jira-api/" name="lib-jira-api" lastBuildLabel="5355" lastBuildTime="2014-05-01T17:55:49Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/libs_svnkit/" name="libs_svnkit" lastBuildLabel="11" lastBuildTime="2012-02-21T05:00:08Z" lastBuildStatus="Failure" activity="Sleeping"/><Project webUrl="http://ci.jenkins-ci.org/job/selenium-tests/" name="selenium-tests" lastBuildLabel="11" lastBuildTime="2012-11-14T18:40:47Z" lastBuildStatus="Failure" activity="Sleeping"/></Projects>`

func Test_reads_valid_cctray_feed(t *testing.T) {
	p := &Projects{}
	r := strings.NewReader(validTrayFeed)

	err := ReadTrayFeed(r, p)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	expectLen := 15
	if len(p.Project) != expectLen {
		t.Fatalf("len(p) = %v, want %v", len(p.Project), expectLen)
	}

	expectUrl := "http://ci.jenkins-ci.org/job/core_selenium-test/"
	if p.Project[0].WebUrl != expectUrl {
		t.Fatalf("p.Project[0].WebUrl = %v, want %v", p.Project[0].WebUrl, expectUrl)
	}

	expectTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2012-10-27 01:55:14 +0000 UTC")
	if p.Project[0].LastBuildTime != expectTime {
		t.Fatalf("p.Project[0].LastBuildTime = %v, want %v", p.Project[0].LastBuildTime, expectTime)
	}
}
