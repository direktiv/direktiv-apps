---
{
  "image": "direktiv/ansible",
  "desc": "Performs an 'ansible-playbook' action. Requires a playbook variable, and a private key (PEM) variable."
}
---

# Ansible

This isolate performs an 'ansible-playbook' action. When using this isolate, 2 variables must be provided:

- playbook
  - contains the yaml-encoded contents of an ansible playbook file, to be actioned by 'ansible-playbook'
- privateKey
  - contains the PEM-encoded contents of a private key file required for access to the remote machine

The following code block demonstrates how to include this isolate in a workflow, while passing the aforementioned variables to it.

```yaml
  - id: ansible
    image: direktiv/ansible:v1
    type: reusable
    files:
      - key: playbook.yml
        scope: workflow
        type: plain
      - key: pk.pem
        scope: workflow
```

## Input

The input object accepted by this isolate contains the following fields:

```yaml
input:
  playbook: playbook.yml
  privateKey: pk.pem
  # collections to install from galaxy before running the playbook
  collections: ["devsec.hardening"]
  args:
    - "-i"
    - "192.168.1.123,"
  show: true # prints the playbook
  envs:
    - "ANSIBLE_STDOUT_CALLBACK=default" # if non-JSON output is wanted
```

The environment variables ANSIBLE_CALLBACK_WHITELIST and ANSIBLE_STDOUT_CALLBACK are set to *json* by default.


*Note: the playbook and privateKey input fields should correspond with the variable names provided in the function declaration.*
