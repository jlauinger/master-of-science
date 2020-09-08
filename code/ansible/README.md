# Ansible provisioning code

This Ansible code is used to provision `vm6`, a virtual machine that I use to collect and analyze data about open-source
projects.

The benefit of using a VM for this is that for long-running data collection steps, I do not need to keep my laptop
running. Also, I don't need to transfer a lot of data to share and backup the data.


## Running Ansible

First, install [Ansible](https://docs.ansible.com/ansible/latest/index.html).

Store the correct vault password in a file called `vaultpass`.

Configure `vm6` as an SSH host and put my SSH public key in the `authorized_keys` file for `root` on `vm6`. Then, on
the development host, execute

```shell
ansible-playbook site.yml -i inventory --vault-password-file=vaultpass
```

This will setup the server.


## Directory structure

 - `group_vars/` contains variables that are applicable to groups of hosts, in this case there is only an `all.yml` file
   that applies to all hosts, which here is only the one `vm6` host.
 - `roles/` contains role definitions, which are the actual provisioning tasks.
 - `vaults/` contains encrypted variables.
 - `inventory` is the file defining the `vm6` target server.
 - `ansible.cfg` contains some general settings for the Ansible CLI tool.
 - `site.yml` defines which roles are applied to `vm6`.
 
 
## Roles

This Ansible playbook configures the following roles:

 - `base`: general setup like shells, message of the day, etc.
 - `ssh`: SSH configuration and authorized users
 - `unattended-upgrades`: configures the daemon to automatically apply security updates
 - `nginx`: sets up a webserver
 - `golang`: installs the Go programming language
 - `jupyter`: installs the Jupyter Notebook software for web-based notebook access
 - `classification`: installs my Python Flask based Go snippet labeling tool, contained in this repository within the
   `data-survey/classification/` directory
 - `acquisition`: installs my Go based data acquisition tool to extract data from open-source Go repositories
