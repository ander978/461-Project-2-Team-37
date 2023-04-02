from API import *
from graphQlAPI import *
from npmAPI import *
from urlParser import *
import sys


if __name__ == "__main__":

    token = config('GITHUB_TOKEN')
    logFilePath = config('LOG_FILE')
    with open(logFilePath, "w") as f:
        print("Error Log Now Running\n", file=f)
    #logLevel = config('LOG_LEVEL')
    urlFilePath = sys.argv[1] #.split("/").pop()
    #urlFilePath = "/Users/haani/Documents/Spring 2023/ECE 461/Haani's Fork/ECE-461-Haani-Repository/Url File.txt"

    validUrls = generateValidUrls(urlFilePath, logFilePath)
    generateGraphQLData(validUrls, token, logFilePath)
    write(validUrls, token, logFilePath)