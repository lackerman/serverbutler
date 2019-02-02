# Serverbutler

A little web app that helps manage configuration for a server.

Build the app
```bash
make
```

To pull docker images from your private Gitlab registry, create a new
[Private Access Token](https://gitlab.com/profile/personal_access_tokens), make 
sure you've checked **read_registry**. 

Then you'll need to perform `docker login` first, using your token as the password.
```bash
docker login registry.gitlab.com -u gitlab-ci-token -p 1gMh4p4oZXdKwkUwpcVz
```
> More info, check [Gitlab Container Registry](https://docs.gitlab.com/ce/user/project/container_registry.html)

## Pull the latest image

```bash
docker pull registry.gitlab.com/lackerman/serverbutler:latest
```