from urllib.parse import urlparse

import requests

# import re
# import urlParser
# from decouple import config


def generateQuery(url):
    parseResults = urlparse(url)
    path = parseResults.path
    repoDetails = path.split("/")
    owner = '"' + repoDetails[1] + '"'
    name = '"' + repoDetails[2] + '"'

    graphQlQuery = (
        """
    query {
      repositoryOwner (login: """
        + owner
        + """) {
        repository (name: """
        + name +
        """) {
          forks {
            totalCount
          }
          issues {
            totalCount
          }
          stargazers {
            totalCount
          }
          watchers {
            totalCount
          }
          pullRequests(last: 100, states: CLOSED) {
            totalCount
            edges {
              node {
                reviewDecision
                additions
                deletions
                commits {
                  totalCount
                }
              }
            }
          }
          defaultBranchRef {
            target {
              ... on Commit {
                history {
                  totalCount
                }
              }
            }
          }
        }
      }
    }"""
    )
    # print(graphQlQuery)
    return graphQlQuery


# A simple function to use requests.post to make the API call. Note the json= section.
def run_graphQlQuery(graphQlQuery, token, logPath):
    headers = {"Authorization": "Bearer " + token}
    request = requests.post(
        'https://api.github.com/graphql', json={"query": graphQlQuery}, headers=headers
    )
    if request.status_code == 200:
        return request.json()
    else:
        with open(logPath, "a") as f:
            print(
                "Query failed to run by returning code of {}. {}".format(
                    request.status_code, graphQlQuery
                ),
                file=f
            )
        raise Exception(
            "Query failed to run by returning code of {}. {}".format(
                request.status_code, graphQlQuery
            )
        )


def generateGraphQLData(validUrls, token, logFilePath):
    graphQlData = []
    for url in validUrls:
        query = generateQuery(url)
        result = run_graphQlQuery(query, token, logFilePath)
        graphQlData.append(result["data"]["repositoryOwner"]["repository"])
    # for i in graphQlData: print("Info - {}".format(i))
    with open("outputGraphQl.txt", "w") as f:
        for i in graphQlData:
            print("Info - {}\n".format(i), file=f)
