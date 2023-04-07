# LMCApi Backend with GORM & GIN

### Features

- Secured auth (register, login)
- Easy to add custom Routes

### Get Started

To get started,
clone this repository with `git clone https://github.com/Luca00711/LMCApi.git`
then change the remote name from origin to lmcapi with
`git remote rename origin lmcapi`.

Now set the repository URL for your project with
`git remote add origin YOUR-REPO-URL`.

From here you can now configure your project, let's start with
editing the `README.md` change it to your project description.

You need to copy the `example.env` to `.env`, edit the `.env`
and change the variables to your project data.

Now you should do an initial commit with `git add .`, `git commit -m "Initial commit"`
and `git push --set-upstream origin main`.

To start the Backend, do `go run main.go`.

### Updating
To update the LMCApi you just need to do `git pull lmcapi main --rebase`,  `git pull --rebase` and `git push`.