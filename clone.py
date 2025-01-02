import os
import subprocess

def main():
    github_user = input("Enter your GitHub username: ").strip()
    clone_dir = input("Enter the directory where you want to clone the repos: ").strip()

    os.makedirs(clone_dir, exist_ok=True)
    os.chdir(clone_dir)

    try:
        result = subprocess.run(
            [
                "gh", "repo", "list", github_user, "--limit", "100", "--json",
                "nameWithOwner,isFork,isArchived", "--jq",
                ".[] | select(.isFork == false and .isArchived == false) | .nameWithOwner"
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )

        if result.returncode != 0:
            print("Failed to fetch repositories. Make sure 'gh' is installed and configured.")
            print(result.stderr)
            return

        repos = result.stdout.strip().splitlines()
    except Exception as e:
        print(f"Error fetching repositories: {e}")
        return

    for repo in repos:
        repo = repo.strip()

        repo_name = repo.split("/")[1] if "/" in repo else repo
        repo_dir = repo_name.replace("/", "_")  # Replace '/' with '_' to avoid invalid directory names
        clone_path = os.path.join(clone_dir, repo_dir)

        os.makedirs(clone_path, exist_ok=True)

        clone_command = f"git clone git@github.com:{repo} {clone_path}"
        try:
            result = subprocess.run(
                ["sh", "-c", clone_command],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )

            if result.returncode == 0:
                print(f"{repo} has been cloned into {clone_path}")
            else:
                print(f"Failed to clone {repo} into {clone_path}")
                print(result.stderr)
        except Exception as e:
            print(f"Error cloning {repo}: {e}")

    print(f"All repos have been cloned into {clone_dir}")

if __name__ == "__main__":
    main()