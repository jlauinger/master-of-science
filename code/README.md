# Identification and analysis of unsafe.Pointer usages in open-source Go code: Implementation

## Data Access

Jupyter Notebook Server aus dem TU Netz erreichbar:

[http://vm6.rbg.informatik.tu-darmstadt.de/notebooks](http://vm6.rbg.informatik.tu-darmstadt.de/notebooks)

HTTP Basic Auth:

 - Benutzer: `johannes`
 - Passwort: `Syp9393`
 
Jupyter Server Passwort, falls benötigt: `admin`

Die Daten liegen auf dem Server im Verzeichnis `/root/data`

Code zum Laden der CSV-Daten:

```python
projects_df = pd.read_csv('/root/data/projects.csv',
                         parse_dates=['project_created_at', 'project_last_pushed_at', 'project_updated_at'])

grep_df = pd.read_csv('/root/data/lexical/grep_findings_0_499.csv')
package_df = pd.read_csv('/root/data/packages_0_499.csv')

vet_df = pd.read_csv('/root/data/lexical/vet_findings_0_499.csv')
gosec_df = pd.read_csv('/root/data/lexical/gosec_findings_0_499.csv')
ast_df = pd.read_csv('/root/data/ast/ast_findings_0_499.csv')
function_df = pd.read_csv('/root/data/ast/functions_0_499.csv')
statement_df = pd.read_csv('/root/data/ast/statements_0_499.csv')
```


## Notes

### Project Scraper

This Go program fetches the top 500 most popular Go projects from Github. Using the Github search, these projects
can be found at this search link:

https://github.com/search?l=Go&o=desc&q=Go&s=stars&type=Repositories

The projects contain huge repositories including:

 - Kubernetes
 - Moby (Docker)
 - Hugo (Website tool)
 - Gin (HTTP server framework)
 - frp (reverse proxy to bypass NAT)
 - Gogs (self-hosted Git service)
 - Syncthing (P2P file sharing)
 - etcd (shared key-value storage cluster)
 - traefik (cloud-native edge router)
 - caddy (multi-platform web server)
 - Ethereum (crypto currency)
 - Gitea (self-hosted Git service)
 - InfluxDB (time-series database)
 - Cockroach (cloud SQL database)
 - Mattermost (Slack alternative)
 - Gorm (object-relational mapper)
 - Hashicorp Vault (Ansible secrets manager)


### Github GraphQL API request

Enter this on https://developer.github.com/v4/explorer

```
query {
  search(query:"language:Go", type:REPOSITORY, first:100) {
    repositoryCount
    edges {
      cursor
      node {
        ... on Repository {
          name
          #descriptionHTML
          #stargazers {
          #  totalCount
          #}
        }
      }
    }
  }
}
```

### Use grep to find out how many unsafe.Pointers exist in module dependencies

```
for dir in $(go mod vendor -v 2>&1 | grep -v "#" | sort | uniq); do 
lines=$(rg unsafe.Pointer vendor/$dir | wc -l); 
echo "$lines $dir"; 
done | sort -n | grep -ve "^0 "
```

There are some ways to show the dependency modules:

 - `go mod vendor -v`: stores the required dependency modules in the `vendor` directory, therefore I can analyze the
   exact version that is used by the module. The `-v` switch prints the modules (and thus their filesystem path) to
   stderr.
 - `go mod vendor && find . -type d vendor`: simply checks all directories in the vendor directory. This is very
   inaccurate because not every directory is a package. The first approach prints the logical modules rather than the
   low-level directories and is thus better.
 - `go list -m all`: using the `-m` flag, the list command prints module information. The `all` keyword instructs it
   to recursively show all dependency modules. This command output includes the versions of the modules, which is very
   nice. Some preprocessing using `cut` is necessary to strip the version and including module from the directory path.
   The problem is that this command lists modules that do not exist in the vendor directory, making the grepping throw
   errors. TODO: find out why this is the case.
 - `go mod graph`: This is very similar to the command above, and shares the problem that not all directories exist. 


## License

Copyright (c) 2020 Johannes Lauinger

Licensed under the terms of the <a rel="license" href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GENERAL PUBLIC LICENSE, Version 3</a>.


