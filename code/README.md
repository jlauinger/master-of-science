# Go Dependency Checker Implementation

## Project Scraper

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


## License

Copyright (c) 2020 Johannes Lauinger

Licensed under the terms of the <a rel="license" href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GENERAL PUBLIC LICENSE, Version 3</a>.


