use std::env;
use std::fs;
use std::io::{self, Write};
use std::process::Command;

fn main() -> io::Result<()> {
    // Prompt for GitHub username
    let mut github_user = String::new();
    print!("Enter your GitHub username: ");
    io::stdout().flush()?;
    io::stdin().read_line(&mut github_user)?;
    let github_user = github_user.trim();

    // Prompt for directory for cloning
    let mut clone_dir = String::new();
    print!("Enter the directory where you want to clone the repos: ");
    io::stdout().flush()?;
    io::stdin().read_line(&mut clone_dir)?;
    let clone_dir = clone_dir.trim();

    // Create the directory if not present
    fs::create_dir_all(clone_dir)?;

    // Navigate to the directory
    env::set_current_dir(clone_dir)?;

    // Fetch all repos
    let output = Command::new("gh")
        .arg("repo")
        .arg("list")
        .arg(github_user)
        .arg("--limit")
        .arg("100")
        .arg("--json")
        .arg("nameWithOwner")
        .arg("--jq")
        .arg(".[].nameWithOwner")
        .output()?;

    if !output.status.success() {
        eprintln!("Failed to fetch repositories. Make sure 'gh' is installed and configured.");
        return Ok(());
    }

    let repos = String::from_utf8_lossy(&output.stdout);

    for repo in repos.lines() {
        let repo = repo.trim();
        
        // Extract the repository name (remove the username part)
        let repo_name = repo.split('/').nth(1).unwrap_or(repo);
        let repo_dir = repo_name.replace('/', "_"); // Replace '/' with '_' to avoid invalid directory names
        let clone_path = format!("{}/{}", clone_dir, repo_dir);

        // Create a subdirectory for the repo
        fs::create_dir_all(&clone_path)?;

        // Clone the repository into the subdirectory
        let clone_command = format!("git clone git@github.com:{} {}", repo, clone_path);

        let status = Command::new("sh")
            .arg("-c")
            .arg(&clone_command)
            .status()?;

        if status.success() {
            println!("{} has been cloned into {}", repo, clone_path);
        } else {
            eprintln!("Failed to clone {} into {}", repo, clone_path);
        }
    }

    println!("All repos have been cloned into {}", clone_dir);

    Ok(())
}

