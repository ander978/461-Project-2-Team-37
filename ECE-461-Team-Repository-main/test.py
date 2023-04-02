try:
    with open('out.txt', 'r') as f:
        lines = f.readlines()
        if len(lines) < 10:
            print("Missing information")
            exit()
        url = lines[0].strip()
        number_of_events = int(lines[1].split(": ")[1])
        number_of_starred = int(lines[2].split(": ")[1])
        number_of_subscribers = int(lines[3].split(": ")[1])
        number_of_commits = int(lines[4].split(": ")[1])
        number_of_open_issues = int(lines[5].split(": ")[1])
        number_of_closed_issues = int(lines[6].split(": ")[1])
        license = lines[7].split(": ")[1]
        community_metric = float(lines[8].split(": ")[1])
        pull_requests = int(lines[9].split(": ")[1])
except FileNotFoundError:
    print("File not found.")
    exit()
except Exception as e:
    print(f"An error occurred while reading the file: {e}")
    exit()

cases_passed = 0
for i in range(1, len(lines)):
    try:
        assert int(lines[i].split(": ")[1]) != 0
        cases_passed += 1
    except:
        pass

#if len(lines) - 1 != cases_passed:
#    print("File format is incorrect. Not all lines are correct")
#    exit()

coverage = (cases_passed / (len(lines) - 1)) * 100
print(f"Total: {len(lines) - 1}")
print(f"Passed: {cases_passed}")
print(f"Coverage: {coverage:.2f}%")
print(f"{cases_passed}/{len(lines) - 1} passed")
