# SRE Operator

Automate repetitive actions using kubernetes.

## Overview

### Triggers

There are 2 ways to trigger an action:

#### Webhooks

You can define a webhook resource:

```yaml
apiVersion: sre.henrywhitaker.com/v1alpha1
kind: Webhook
metadata:
    name: demo-webhook
spec:
    id: demo-hook
```

Then you can trigger it manually by running:

```
curl -X POST {domain}/webhook/demo-hook
```

Or setup your monitoring/alerting platform to call the webhook when a monitor is triggered.

#### Schedule

You can schedule actions to be triggered based on a cron schedule:

```yaml
apiVersion: sre.henrywhitaker.com/v1alpha1
kind: Schedule
metadata:
    name: demo-schedule
spec:
    id: demo-cron
    cron: 0 * * * *
```

### Actions

You can then create actions that subscribe to defined triggers:

#### Rollout

You can trigger the equivalent of `kubectl rollout {action} {kind} {name}`:

```yaml
apiVersion: sre.henrywhitaker.com/v1alpha1
kind: Rollout
metadata:
    name: demo-rollout
spec:
    triggers:
        - demo-hook
        - demo-cron
    target:
        kind: deployment
        name: coredns
        namespace: kube-system
    # The rolout action, one of: restart, pause, resume
    action: restart
    # This is an optional field to throttle actions once for
    # every specified period
    throttle: 10m
```

#### Script

You can define a script to be run:

```yaml
apiVerison: sre.henrywhitaker.com/v1alpha1
kind: Script
metadata:
    name: demo-script
spec:
    triggers:
        - demo-hook
        - demo-cron
    # Optional: specifcy the image to use, defaults to alpine:latest
    image: alpine:latest
    # Optional: specify the shell to run the script with, defaults to: /bin/sh
    shell: /bin/sh
    script: |
        echo Demo script!
    # Optional: define secrets to be mounted to /var/sre/secrets and to populate env with
    # These must be in the same namespace as the sre-operator
    secrets:
        - example-secret
```
