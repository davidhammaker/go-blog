# GO Blog

_A minimal blog built with GO._

## Getting Started

### Create a Docker Network

```sh
docker network create go-blog-network
```

### Run the Blog

#### Build the container:

```sh
docker build -t go-blog-container .
```

#### Run:

Set the following environment variables:

| Variable              | Description                                                                                                                                               |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `MYSQL_ROOT_PASSWORD` | The root password you're using for your MySQL database.                                                                                                   |
| `DBUSER`              | If not `"root"`, you may need to configure some permissions.                                                                                              |
| `DBPASS`              | The password for `DBUSER`; if the user is `"root"`, this will match `MYSQL_ROOT_PASSWORD`.                                                                |
| `DBHOST`              | For development, this will be the name of your MySQL container. E.g. `"go-blog-mysql"`.                                                                   |
| `DBPORT`              | Most likely 3306, unless you choose to change it.                                                                                                         |
| `DBNAME`              | Any database name you choose. E.g. `"mydbname"`.                                                                                                          |
| `HOMEREF`             | The URL reference to a markdown file representing the blog's home page. E.g. `"https://raw.githubusercontent.com/davidhammaker/go-blog/master/README.md"` |
| `BLOGTITLE`           | The title of your blog, as it will appear on the top of each page.                                                                                        |

<br />

Then run:

```sh
docker run -d --name go-blog-app --network go-blog-network -p 8080:8080 \
  -e DBUSER=$DBUSER \
  -e DBPASS=$DBPASS \
  -e DBHOST=$DBHOST \
  -e DBPORT=$DBPORT \
  -e DBNAME=$DBNAME \
  -e HOMEREF=$HOMEREF \
   -e BLOGTITLE="$BLOGTITLE" \
  go-blog-container
```

#### Stop or Restart the container:

```sh
docker stop go-blog-app
docker restart go-blog-app
```

### Run MySQL

#### Build:

```sh
docker build -t go-blog-mysql -f mysql.Dockerfile .
```

#### Run:

Set the `MYSQL_ROOT_PASSWORD` (same as you set for running the container), and then run:

```sh
docker run -v .:/usr/src --name go-blog-mysql --network go-blog-network -p 3306:3306 -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -d go-blog-mysql
```

#### Create database:

```sh
docker exec -it go-blog-mysql /bin/bash
```

_Then, from the `bash` prompt:_

```sh
mysql -p
```

_Then, from the `mysql` prompt:_

```
create database mydbname
```

#### Upgrade:

```sh
docker exec -it go-blog-mysql /bin/bash
```

_Then run the following commands in the `bash` prompt:_

```sh
mysql -p -D mydbname < migrations/upgrade/0000_migration_table.sql
mysql -p -D mydbname < migrations/upgrade/0001_initial_migration.sql
```

- Downgrading is similar.

#### Stop and Restart:

```sh
docker stop go-blog-mysql
docker restart go-blog-mysql
```
