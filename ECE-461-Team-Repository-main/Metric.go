package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

func readFromFile(file string) string {
	data := ""

	f, err := os.Open(file)
	if err != nil {
		// fmt.Println(err)
		return data
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data += scanner.Text() + "\n"
	}

	return data
}

func getDependencyInfo(dependenciesString string) float64 {
	regex := regexp.MustCompile(`([\w-]+):\s*([^,]+)`)
	dependencies := regex.FindAllStringSubmatch(dependenciesString, -1)
	numDependencies := len(dependencies)
	if numDependencies == 0 {
		return 1.0
	}
	numPinnedDependencies := 0
	for _, dependency := range dependencies {
		_, version := dependency[1], dependency[2]
		version = strings.TrimSpace(version)
		match, _ := regexp.MatchString(`^\^0\.\d+(\.\d+)*$`, version)
		if match {
			numPinnedDependencies++
			continue
		}
		match, _ = regexp.MatchString(`^~\d+\.\d+(\.\d+)*$`, version)
		if match {
			numPinnedDependencies++
			continue
		}
		match, _ = regexp.MatchString(`^\d+\.\d+(\.\d+)*(\.[xX*])?$`, version)
		if match {
			numPinnedDependencies++
			continue
		}
		match, _ = regexp.MatchString(`^\d+\.\d+(\.\d+)*$`, version)
		if match {
			numPinnedDependencies++
			continue
		}
		match, _ = regexp.MatchString(`^\d+\.\d+(\.\d+)*\s-\s\d+\.\d+(\.\d+)*$`, version)
		if match {
			numPinnedDependencies++
		}
	}

	return float64(numPinnedDependencies) / float64(numDependencies)
}



// New metric - Code Reviewed PRs
func prCodeReviewScore(totalCommits int, totalPRs int, prCounts [6]int) float64 {
    // Keep everything float for convenience
    var crRecency float64
    var crRecencyScore float64
    var crprRatio float64
    var commitAvg float64
    var estGoodCommits float64
    var commitRatio float64
    var commitRatioScore float64
    const (
	    good_pr_ind = 0
	    bad_pr_ind = 1
	    good_line_ind = 2
	    bad_line_ind = 3
	    good_commit_ind = 4
	    bad_commit_ind = 5
	)
    crRecencyWeight := 0.30
    commitRatioWeight := 0.70

    // Type conversion for more accurate scores
    goodPRs := float64(prCounts[good_pr_ind])
    badPRs := float64(prCounts[bad_pr_ind])
    goodLines := float64(prCounts[good_line_ind])
    badLines := float64(prCounts[bad_line_ind])
    goodCommits := float64(prCounts[good_commit_ind])
    badCommits := float64(prCounts[bad_commit_ind])

    // GraphQL got the 100 most recent PRs (sometimes this breaks in the Explorer for some reason)
    // The first fraction is how many of the recently added lines of code come from CRPRs.
    crRecency = goodLines / (badLines + goodLines)
    crRecencyScore = crRecency * crRecencyWeight

    // The second fraction estimates how many of all the commits on the default branch come from CRPRs.
    crprRatio = goodPRs / (badPRs + goodPRs)
    // Technically don't need to split commit count by good and bad, but could be useful?
    commitAvg = ((goodCommits + badCommits) / 100) + 1
    // +1 for the merge commit that GraphQL does not count.
    estGoodCommits = float64(totalPRs) * crprRatio * commitAvg

    commitRatio = estGoodCommits / float64(totalCommits)
    commitRatioScore = commitRatio * commitRatioWeight

    return (crRecencyScore + commitRatioScore)
}

// Use lines of code, as the more lines there are the harder it will be to learn
// Use community metric, as reflects different methods of help access such as readme and license
//
//export rampUpScore
func rampUpScore(communityMetric int, linesOfCode int) float64 {
	metricScale := float64(communityMetric) / 100
	linesScale := float64(linesOfCode) / 5000
	if linesScale > 1 {
		linesScale = 1
	}
	linesScale = 1 - linesScale
	return ((metricScale + linesScale) / 2)
}

//export license
func licenseScore(license string) float64 {
    // Is this the right license we should be looking for?
	if license == "MIT License" {
		return 1.0
	} else {
		return 0.0
	}
}

//export busFactorScore
func busFactorScore(forks int, lines int, pulls int) float64 {
	forksScore := float64(forks) / 500
	if forksScore > 1 {
		forksScore = 1
	}
	linesScale := float64(lines) / 5000
	if linesScale > 1 {
		linesScale = 1
	}
	linesScale = 1 - linesScale
	pullsScore := float64(pulls) / 500
	if pullsScore > 1 {
		pullsScore = 1
	}
	score := (linesScale + forksScore + pullsScore) / 3
	return float64(score)
}

//export correctnessScore
func correctnessScore(openIssues int, closedIssues int, starred int, subscribers int) float64 {
	//totalIssues := openIssues + closedIssues
	openIssuesRatio := 1 - (float64(openIssues) / float64(closedIssues))
	if openIssuesRatio > 1 {
		openIssuesRatio = 1
	}
	subscribersScore := float64(subscribers) / 100
	starredScore := float64(starred) / 500
	if starredScore > 1 {
		starredScore = 1
	}
	score := (openIssuesRatio + subscribersScore + starredScore) / 3
	return score
}

func responsiveMaintainerScore(commits int, closedIssues int) float64 {
	commitsScore := float64(commits) / 100
	closedIssuesScore := float64(commits) / 100
	return (commitsScore + closedIssuesScore) / 2.0
}

// New metrics - need to determine weighting
func netScore(rampUp float64,
              correctness float64,
              busFactor float64,
              responsiveMaintainer float64,
              license float64,
              prCodeReview float64,
	      dependencies float64,
              ) float64 {
	final_score := (2 * rampUp +
	    2   * correctness +
	    1.5 * busFactor +
	    3   * responsiveMaintainer +
	    2   * license +
	    2   * prCodeReview +
	    1.5 * dependencies ) / 12.5
	return final_score
}

func numLines(dir string) int {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return 2500
	}
	out := 0
	for _, f := range files {
		//fmt.Println(f.Name())
		content, err := ioutil.ReadFile(dir + "/" + f.Name())
		content1 := string(content)
		if err != nil {
			//fmt.Println("Error reading file:", err)
			continue
		}
		lines := strings.Split(content1, "\n")
		nonEmpty := []string{}
		for _, str := range lines {
			ex := ([]rune(str))
			ex1 := 13
			if len(ex) != 0 {
				ex1 = int(ex[0])
			}
			if ex1 != 13 || len(ex) != 1 {
				nonEmpty = append(nonEmpty, str)
			}
		}
		out = out + len(nonEmpty)
		// fmt.Println(len(nonEmpty))
	}
	return out
}

type pair struct {
	Key   string
	Value float64
}

func main() {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./github_scores.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the scores table
	sqlStmt := `
	create table if not exists scores (
		id integer not null primary key,
		url text,
		package_name text,
		net_score float,
		ramp_up_score float,
		correctness_score float,
		bus_factor_score float,
		responsive_maintainer_score float,
		license_score float,
		pr_review_score float,
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	// Load in REST file
	file := "out.txt"
	data := readFromFile(file)
	lines := strings.Split(data, "https")
	lines = lines[1:]

	// Load in graphQL file
	file_QL := "outputGraphQl.txt"
	data_QL := readFromFile(file_QL)
	lines_QL := strings.Split(data_QL, "}}}}")
	lines_QL = lines_QL[:(len(lines_QL) - 1)]

	// url and packageName init
	var url string
	var packageName string

	// Declare the name of variable to use for later calculations
	var Number_of_Events int
	var Number_of_Starred int     //
	var Number_of_Subscribers int //
	var Number_of_Commits int     //
	var Number_of_Open_Issues int
	var Number_of_Closed_Issues int //
	var Community_Metric int        //
	var Pull_Requests int           //
	var Number_of_Watchers int
	var Lines_of_Code int   //*
	var Number_of_Forks int //
	var Number_of_Total_Issues int
	var License string
	var Dependencies string

	// Variable names for pr_review_score
	var Number_of_Total_Commits int
	var Review_Result string
	var Changed_Lines int
	var PR_Commit_Count int
	var Is_Code_Reviewed bool
	// Storage of multiple counts for pr_review_score
	var PR_Review_Counts = [6]int{0, 0, 0, 0, 0, 0}
	const (
	    good_pr_ind = 0
	    bad_pr_ind = 1
	    good_line_ind = 2
	    bad_line_ind = 3
	    good_commit_ind = 4
	    bad_commit_ind = 5
	)

	scores := make(map[string]float64)
	for i, line := range lines {
		line1 := strings.Split(line, "\n")
		line1[0] = "https" + line1[0]

		// Get url and package name
		noApiParts := strings.Split(line1[0], "api.")
	    noRepoParts := strings.Split(noApiParts[0] + noApiParts[1], "repos/")
	    url = noRepoParts[0] + noRepoParts[1]
        splitUrl:= strings.Split(url, "/")
        packageName = splitUrl[(len(splitUrl)-1)]

		// Get line count
		dir := strings.Split(line1[0], "/") //[len(line1)-1]
		dir1 := dir[len(dir)-1]
		Lines_of_Code = numLines(dir1)

		// Parse out.txt
		for _, ind := range line1 {
			if strings.Contains(ind, "Number of Events") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[3], "%d", &Number_of_Events)
			} else if strings.Contains(ind, "Number of Subscribers") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[3], "%d", &Number_of_Subscribers)
			} else if strings.Contains(ind, "Number of Commits") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[3], "%d", &Number_of_Commits)
			} else if strings.Contains(ind, "Number of Open_Issues") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[3], "%d", &Number_of_Open_Issues)
			} else if strings.Contains(ind, "Number of Closed_Issues") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[3], "%d", &Number_of_Closed_Issues)
			} else if strings.Contains(ind, "Community Metric") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2], "%d", &Community_Metric)
			} else if strings.Contains(ind, "License") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[1], "%s", &License)
			} else if strings.Contains(ind, "Dependencies") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[1], "%s",&Dependencies) 
			}
		}
		// Reset counts for every query
		PR_Review_Counts = [6]int{0, 0, 0, 0, 0, 0}

        // Parse outputGraphQl.txt
		linesQL1 := strings.Split(lines_QL[i], ",")
		for _, ind := range linesQL1 {
			if strings.Contains(ind, "forks") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[4][:(len(fields[4])-1)], "%d", &Number_of_Forks)
			} else if strings.Contains(ind, "issues") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2][:(len(fields[2])-1)], "%d", &Number_of_Total_Issues)
			} else if strings.Contains(ind, "stargazers") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2][:(len(fields[2])-1)], "%d", &Number_of_Starred)
			} else if strings.Contains(ind, "watchers") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2][:(len(fields[2])-1)], "%d", &Number_of_Watchers)
			} else if strings.Contains(ind, "pullRequests") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2][:(len(fields[2]))], "%d", &Pull_Requests)
			} else if strings.Contains(ind, "defaultBranchRef") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[4][:(len(fields[4]))], "%d", &Number_of_Total_Commits)
			} else if strings.Contains(ind, "edges") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[3][1:(len(fields[3]))-1], "%s", &Review_Result)
				if Review_Result == "null" {
					Is_Code_Reviewed = false
					PR_Review_Counts[bad_pr_ind] += 1
				} else {
					Is_Code_Reviewed = true
					PR_Review_Counts[good_pr_ind] += 1
				}
			} else if strings.Contains(ind, "node") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2][1:(len(fields[2]))-1], "%s", &Review_Result)
				if Review_Result == "null" {
					Is_Code_Reviewed = false
					PR_Review_Counts[bad_pr_ind] += 1
				} else {
					Is_Code_Reviewed = true
					PR_Review_Counts[good_pr_ind] += 1
				}
			} else if strings.Contains(ind, "additions") || strings.Contains(ind, "deletions") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[1][:(len(fields[1]))], "%d", &Changed_Lines)
				if Is_Code_Reviewed {
					PR_Review_Counts[good_line_ind] += Changed_Lines
				} else {
					PR_Review_Counts[bad_line_ind] += Changed_Lines
				}
			} else if strings.Contains(ind, "commits") {
				fields := strings.Fields(ind)
				fmt.Sscanf(fields[2][:(len(fields[2])-3)], "%d", &PR_Commit_Count)
				if Is_Code_Reviewed {
					PR_Review_Counts[good_commit_ind] += PR_Commit_Count
				} else {
					PR_Review_Counts[bad_commit_ind] += PR_Commit_Count
				}
			}
		}

		scores["RAMP_UP_SCORE"] = rampUpScore(Community_Metric, Lines_of_Code)
		scores["CORRECTNESS_SCORE"] = correctnessScore(Number_of_Open_Issues, Number_of_Closed_Issues, Number_of_Starred, Number_of_Subscribers)
		scores["BUS_FACTOR_SCORE"] = busFactorScore(Number_of_Forks, Lines_of_Code, Pull_Requests)
		scores["RESPONSIVE_MAINTAINER_SCORE"] = responsiveMaintainerScore(Number_of_Commits, Number_of_Closed_Issues)
		scores["LICENSE_SCORE"] = licenseScore(License)
		// New metrics - scores
		scores["DEPENDCIES"] = getDependencyInfo(Dependencies)
		scores["PR_REVIEW_SCORE"] = prCodeReviewScore(Number_of_Total_Commits, Pull_Requests, PR_Review_Counts)
		net_score := netScore(scores["RAMP_UP_SCORE"],
		                      scores["CORRECTNESS_SCORE"],
		                      scores["BUS_FACTOR_SCORE"],
		                      scores["RESPONSIVE_MAINTAINER_SCORE"],
		                      scores["LICENSE_SCORE"],
				      scores["PR_REVIEW_SCORE"],scores["DEPENDCIES"])
		keys := make([]pair, 0, len(scores))
		for key, value := range scores {
			keys = append(keys, pair{key, value})
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Value > keys[j].Value
		})

		// Insert the GitHub URL into the scores table
		// This feels like test code, idk if it was needed
		// I redefined url and packageName above

        // url := "https://github.com/my_username/my_repository"
        // packageName := "my_package_name"
        // tx, err := db.Begin()
        // if err != nil {
        //     log.Fatal(err)
        // }
        // stmt, err := tx.Prepare("insert into scores(url, package_name) values(?, ?)")
        // if err != nil {
        //     log.Fatal(err)
        // }
        // defer stmt.Close()
        // _, err = stmt.Exec(url, packageName)
        // if err != nil {
        //     log.Fatal(err)
        // }
        // tx.Commit()
        //
        // fmt.Println("GitHub URL stored in the database")

        // NPM ingestibility: a package is only uploaded if all part 1 scores >= 0.5
        if (scores["RAMP_UP_SCORE"] >= 0.5 &&
        scores["CORRECTNESS_SCORE"] >= 0.5 &&
        scores["BUS_FACTOR_SCORE"] >= 0.5 &&
        scores["RESPONSIVE_MAINTAINER_SCORE"] >= 0.5 &&
        scores["LICENSE_SCORE"] >= 0.5) {

            // Insert the scores into the database
            tx, beginErr := db.Begin()
            if beginErr != nil {
                log.Fatal(beginErr)
            }
            stmt, prepareErr := tx.Prepare(`
                insert into scores(
                    url, package_name, net_score,
                    ramp_up_score, correctness_score, bus_factor_score,
                    responsive_maintainer_score, license_score, pr_review_score
                ) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            `)
            if prepareErr != nil {
                log.Fatal(prepareErr)
            }
            defer stmt.Close()
            _, err = stmt.Exec(
                url, packageName, net_score,
                scores["RAMP_UP_SCORE"], scores["CORRECTNESS_SCORE"], scores["BUS_FACTOR_SCORE"],
                scores["RESPONSIVE_MAINTAINER_SCORE"], scores["LICENSE_SCORE"], scores["PR_REVIEW_SCORE"],
            )
            if err != nil {
                log.Fatal(err)
            }
            err = tx.Commit()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Scores added to the database.")
            // keys = sort.Sort(sort.Reverse(sort.StringSlice(keys)))
            // Do we still need this next line?
            fmt.Printf("{\"URL\":\"%s\", \"NET_SCORE\":%0.2f, \"%s\":%0.2f, \"%s\":%0.2f, \"%s\":%0.2f, \"%s\":%0.2f}\n",
                url, net_score, keys[0].Key, scores[keys[0].Key], keys[1].Key, scores[keys[1].Key], keys[2].Key, scores[keys[2].Key], keys[3].Key, scores[keys[3].Key])
        }
    }
}

















