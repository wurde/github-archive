# GitHub Archive

Backup GitHub repositories with ease.

## Usage

Start by generating a GitHub Personal access token:

1. Ensure you've verified your email address.
2. Click your profile photo, then click Settings.
3. In the left sidebar, click Developer settings.
4. In the left sidebar, click Personal access tokens.
5. Click Generate new token.
6. Give your token a descriptive name. 
7. Select the Expiration via drop-down menu.
8. Select the scopes, or permissions to grant.
9. Click Generate token.
10. Set the variable `GITHUB_TOKEN=` in `./.env`.

        Warning: Treat your tokens like passwords
        and keep them secret. When working with the
        API, use tokens as environment variables
        instead of hardcoding them into your programs.

Update a `./repos.txt` file with target repos:

```txt
github.com/example/dotfiles
github.com/example/personal-site
github.com/example/scripts
```

Run the program:

```bash
git clone https://github.com/wurde/github-archive
cd github-archive
./github-archive
```
