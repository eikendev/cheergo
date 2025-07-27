<div align="center">
	<h1>cheergo</h1>
	<h4 align="center">
		Stay in the loop when your projects get noticed.
	</h4>
	<p>
		<strong>cheergo</strong> notifies you when your GitHub repositories are starred or followed, so you never miss a moment of recognition.
	</p>
</div>

<p align="center">
	<a href="https://github.com/eikendev/cheergo/actions"><img alt="Build status" src="https://img.shields.io/github/actions/workflow/status/eikendev/cheergo/main.yml?branch=main"/></a>&nbsp;
	<a href="https://github.com/eikendev/cheergo/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/eikendev/cheergo"/></a>&nbsp;
</p>

## ✨ Why cheergo?

Ever wondered who’s cheering for your open source work? **cheergo** keeps you connected to your community by sending you notifications whenever someone stars or follows your repositories. Whether you’re a solo dev or part of a team, cheergo helps you celebrate every milestone.

## 🚀 Features

- **Flexible delivery**: Send alerts to email, Telegram, Slack, Discord, and more (powered by [Shoutrrr](https://containrrr.dev/shoutrrr/latest/services/overview/))
- **Smart summaries**: Get concise, AI-generated notifications with OpenAI/OpenRouter (optional)
- **Easy setup**: Configure via CLI flags or environment variables

## 🛠️ How It Works

1. **cheergo** checks your GitHub account for new stars and followers.
2. It compares the latest state with your previous data (stored locally in a YAML file).
3. When it detects something new, it crafts a notification, optionally using an LLM for a smart summary.
4. The message is sent to your chosen channel(s) via Shoutrrr.

## 📦 Installation

**Recommended:** Download the latest binary from the [releases page](https://github.com/eikendev/cheergo/releases).

**Or build from source:**
```bash
go install github.com/eikendev/cheergo/cmd/...@latest
```

## ⚡ Quick Start

1. **Set up your notification channel** (see [Shoutrrr docs](https://containrrr.dev/shoutrrr/latest/services/overview/) for supported services).

2. **Run cheergo** with your GitHub username and Shoutrrr URL:

```bash
cheergo run --github-user YOUR_GITHUB_USERNAME --shoutrrr-url YOUR_SHOUTRRR_URL
```

3. To enable **AI-powered summaries**, add your OpenAI/OpenRouter API key:

```bash
cheergo run --github-user ... --shoutrrr-url ... --llm-api-key YOUR_API_KEY
```

_All options can also be set via environment variables (see below)._

## ⚙️ Configuration

| Option            | CLI Flag             | Env Variable              | Default                | Description                                 |
|-------------------|---------------------|---------------------------|------------------------|---------------------------------------------|
| Storage file      | `--storage`         | `CHEERGO_STORAGE`         | `storage.yml`          | Path to local storage file                  |
| Shoutrrr URL      | `--shoutrrr-url`    | `CHEERGO_SHOUTRRR_URL`    | *(required)*           | Notification channel URL                    |
| GitHub user       | `--github-user`     | `CHEERGO_GITHUB_USER`     | *(required)*           | GitHub username to monitor                  |
| LLM API key       | `--llm-api-key`     | `CHEERGO_LLM_API_KEY`     | *(optional)*           | OpenAI/OpenRouter API key for summaries     |
| LLM base URL      | `--llm-base-url`    | `CHEERGO_LLM_BASE_URL`    | `https://openrouter.ai/api/v1` | LLM API endpoint                |
| LLM model         | `--llm-model`       | `CHEERGO_LLM_MODEL`       | `google/gemini-2.5-flash-lite-preview-06-17` | LLM model to use         |
| Verbose logging   | `--verbose`         | `CHEERGO_VERBOSE`         | `false`                | Enable debug logging                        |

## 🔔 Supported Notification Channels

cheergo uses [Shoutrrr](https://containrrr.dev/shoutrrr/latest/services/overview/) for notifications, supporting:
- Email (SMTP)
- Telegram
- Slack
- Discord
- Microsoft Teams
- Matrix
- Rocket.Chat
- ...and many more!

Just provide the appropriate Shoutrrr URL for your service.

## 🧩 Libraries Used

- [`kong`](https://github.com/alecthomas/kong) – CLI parsing
- [`shoutrrr`](https://github.com/containrrr/shoutrrr) – Multi-channel notifications
- [`go-github`](https://github.com/google/go-github) – GitHub API
- [`go-openai`](https://github.com/sashabaranov/go-openai) – LLM integration (optional)
- [`yaml.v3`](https://pkg.go.dev/gopkg.in/yaml.v3) – YAML config/state
