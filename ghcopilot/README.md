# GitHub Copilot Dagger Module

An example repo of using GitHub Copilot CLI in a dagger module.

Similar to other coding agents, we want to show how to take advantage of GitHub Copilot CLI to multi-step or parallel agentic workloads.


In bash ensure you export the ```GITHUB_TOKEN``` environment variable which holds your GH PAT Token with ```Copilot Requests``` permission added and set to ```Allow Read Only```
```bash
export GITHUB_TOKEN="gh_pat..."
```

You can then run the dagger shell and run accordingly
```bash
dagger

# Dagger Shell
new-ghcopilot env://GITHUB_TOKEN | with-prompt "some inline prompt" | response
```

Or System Shell
```bash
dagger -c "new-ghcopilot env://GITHUB_TOKEN | with-prompt 'some inline prompt' | response"
```

Or with the Dagger CLI
```bash
dagger call new-ghcopilot --token=env://GITHUB_TOKEN with-prompt --prompt='some inline prompt' response
```

## Notes
Q: Why must I use a GH PAT token, can't I just pass in my GitHub CLI OAuth Token (```gh auth login```)??
A: The GitHub CLI has long since implemented using the local keyring solution on a given OS.  As such when you run ```gh auth login``` the command will securely store your token into the local OS keyring instead of a plain text value stored in ~/.config/gh/hosts.yml.  You can sorta follow issues/concerns about the previous insecure way of doing this [here](https://github.com/cli/cli/issues/8954) - I'm sure there are more detailed discussions/articles/issues about this as well.