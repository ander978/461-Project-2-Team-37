import validators
from npmAPI import findGitUrl
import git

#urlFilePath = "/Users/haani/Documents/Spring 2023/ECE 461/Haani's Fork/ECE-461-Haani-Repository/Url File.txt"

def clone(url):
    s = url.split("/").pop()
    local_path = s
    #os.makedirs(s)
    git.Repo.clone_from(url, local_path)


def parseUrls(urlFilePath, logFilePath):
    try:
        fp = open(urlFilePath)
        urls = fp.readlines()
    except:
        with open(logFilePath, "a") as f:
            print("URL Text file either does not exist or is in incorrect directory", file=f)
    finally:
        fp.close()
    return urls

def urlValidator(urls, logFilePath):
    validatedUrls = []
    for url in urls:
        url = url.replace('\n','')
        if validators.url(url) is True and 'github.com' in url:
            try:
                clone(url)
            except:
                # print('Cannot Clone. No Access to or already cloned', url)
                #os.environ["LOG_LEVEL"] = 2
                with open(logFilePath, "a") as f:
                    print("Cannot Clone. No Access to or already cloned {}\n".format(url), file=f)
                # log that repo could not be cloned
            validatedUrls.append(url)
        elif validators.url(url) is True and 'www.npmjs.com' in url:
            gitUrl = findGitUrl(url)
            try:
                clone(gitUrl)
            except:
                # print('Cannot Clone. No Access to or already cloned', url)
                with open(logFilePath, "a") as f:
                    print("Cannot Clone. No Access to or already cloned {}\n".format(url), file=f)
                 # Log that repo could not be cloned
            validatedUrls.append(gitUrl)
        else:
            print('The URL could not be validated: ', url)
            with open(logFilePath, "a") as f:
                print("The URL could not be validated {}\n".format(url), file=f)

    for i in range(len(validatedUrls)):
        if 'ssh://git@' in validatedUrls[i]:
            validatedUrls[i] = validatedUrls[i].replace('ssh://git@','https://')
    return validatedUrls

def generateValidUrls(urlFilePath, logFilePath):
    urls = parseUrls(urlFilePath, logFilePath)
    validUrls = urlValidator(urls, logFilePath)
    return validUrls
