# Contributions farmer

**⚡ Blazingly fast contributions farmer. </br> 🔒 Use it only if you want to flex about your contributions count to your parents.**

![Farmer's result](.github/assets/result.png)

I created this project because of the speed of other scripts (they are not as fast as mine 😎). </br>
You can deploy this script on Railway or another infrastructure platform, but I assure you it is not necessary because of its speed.</br>

## Step-by-step

**First of all you need to generate two github tokens.**

- [Fine-grained access token](https://github.com/settings/tokens?type=beta)
- [Classic token](https://github.com/settings/tokens)

**Create your own repository and clone this app**

```bash
git clone https://github.com/robotiksuperb/contributions-farmer.git
```

**Configure `app.env` file.** </br>
**Farmer will not work if you don't configure these fields:**

```
ACCESS_TOKEN=
CLASSIC_TOKEN=
USER_NAME=
USER_EMAIL=
```

It's not neccessarry to configure other parameters, but it's on your own.

And the last one, you can easily change the dates for the farmer in the `main.go` file.

## How-to-run

- Windows

```bash
make all-w
```

- Linux

```bash
make all-l
```
