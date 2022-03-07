module netflix_agent

go 1.14

replace netflix_agent/utils => ./utils

require (
	github.com/evsio0n/log v0.0.0-20210802084243-341ca719db1b
	netflix_agent/utils v0.0.0-00010101000000-000000000000
)
