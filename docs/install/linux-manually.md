## Manual installation and configuration

### Install

Simply download one of the binaries for your system:

```bash
sudo wget -O /usr/local/bin/gitlab-ci-multi-runner https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-386
sudo wget -O /usr/local/bin/gitlab-ci-multi-runner https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-amd64
sudo wget -O /usr/local/bin/gitlab-ci-multi-runner https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-arm
```

Give it permissions to execute:

```bash
sudo chmod +x /usr/local/bin/gitlab-ci-multi-runner
```

Optionally, if you want to use Docker, install Docker with:

```bash
curl -sSL https://get.docker.com/ | sh
```

Create a GitLab CI user (on Linux):

```
sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash
```

Register the runner:

```bash
sudo gitlab-ci-multi-runner register
```

Install and run as service (on Linux):
```bash
sudo gitlab-ci-multi-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
sudo gitlab-ci-multi-runner start
```

### Update

Stop the service (you need elevated command prompt as before):

```bash
sudo gitlab-ci-multi-runner stop
```

Download the binary to replace runner's executable:

```bash
sudo wget -O /usr/local/bin/gitlab-ci-multi-runner https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-386
sudo wget -O /usr/local/bin/gitlab-ci-multi-runner https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-amd64
```

Give it permissions to execute:

```bash
sudo chmod +x /usr/local/bin/gitlab-ci-multi-runner
```

Start the service:

```bash
sudo gitlab-ci-multi-runner start
```
