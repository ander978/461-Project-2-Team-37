import requests
import validators

urlFilePath = "/Users/haani/Documents/Spring 2023/ECE 461/Haani's Fork/ECE-461-Haani-Repository/Url File.txt"

def findGitUrl(url):
    #print('Function code 1', url)
    splits = url.split('/')
    packageName = splits[-1]
    url = 'https://registry.npmjs.org/' + packageName
    #print('Function code 2', url)
    try:
        response = requests.get(url)
    except:
        print('Cannot get response from NPM API')
        with open(logFilePath, "a") as f:
            print("Cannot get response from NPM API {}\n".format(response), file=f)
    #print('Function code 3', response)
    rd = response.json()
    gitUrl = rd['repository']['url']
    gitUrl = gitUrl[4:]
    gitUrl = gitUrl.replace('.git', '')
    #print('Succes. Git URL is', gitUrl)

    return gitUrl

#findGitUrl('https://www.npmjs.com/package/express')