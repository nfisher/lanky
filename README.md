# Lanky

[![Build Status](https://travis-ci.org/nfisher/lanky.svg)](https://travis-ci.org/nfisher/lanky)

The thinner faster younger brother of Janky. The primary goal of this version is to make it as stateless as possible while minimising the dependencies on the server its deployed to as well.

A fundamental deviation is that projects registered in jenkins will be tied back to GitHub in the following way;

Jenkins Project Name: 

```
${GITHUB_REPOSITORY_NAME}-${GITHUB_REPOSITORY_ID}
```

Request payload signatures are verified using the HMAC signing that uses the GitHub secret.
