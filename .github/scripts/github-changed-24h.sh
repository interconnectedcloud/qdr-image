ORG="${ORG:-apache}"
REPO="${REPO:-qpid-dispatch}"

COMMITS_JSON=`curl -s https://api.github.com/repos/${ORG}/${REPO}/commits`
LAST_COMMIT=`echo "${COMMITS_JSON}" | jq -r '.[0].commit.committer.date'`
SHA=`echo "${COMMITS_JSON}" | jq -r '.[0].sha'`
MSG=`echo "${COMMITS_JSON}" | jq -r '.[0].commit.message'`
AUTHOR=`echo "${COMMITS_JSON}" | jq -r '.[0].commit.author.email'`

LAST_COMMIT_TIMESTAMP=`date -u -d "${LAST_COMMIT}" +%s`
NOW_TIMESTAMP=`date -u +%s`
ELAPSED_SECS=$((NOW_TIMESTAMP - LAST_COMMIT_TIMESTAMP))

if [[ ${ELAPSED_SECS} -lt 86400 ]]; then
    echo "Last commit happened in the past 24h: ${LAST_COMMIT}"
    echo "::set-output name=changed::true"
else
    echo "Last commit is past 24h: ${LAST_COMMIT}"
    echo "::set-output name=changed::false"
fi
echo "${SHA} - ${AUTHOR} - ${MSG}"
