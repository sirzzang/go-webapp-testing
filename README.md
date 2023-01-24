- database 폴더가 이미 초기화되어 있으므로, sql문 동작하지 않음
```bash
 ~/tutorial/go-webapp | master ?3  docker-compose up                        ok | 23:22:26 
Creating network "go-webapp_default" with the default driver
Creating go-webapp_postgres_1 ... done
Attaching to go-webapp_postgres_1
postgres_1  | 
postgres_1  | PostgreSQL Database directory appears to contain a database; Skipping initialization
postgres_1  | 
postgres_1  | 2023-01-24 14:22:30.666 UTC [1] LOG:  starting PostgreSQL 14.5 (Debian 14.5-2.pgdg110+2) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
postgres_1  | 2023-01-24 14:22:30.666 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
postgres_1  | 2023-01-24 14:22:30.666 UTC [1] LOG:  listening on IPv6 address "::", port 5432
postgres_1  | 2023-01-24 14:22:30.672 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
postgres_1  | 2023-01-24 14:22:30.698 UTC [28] LOG:  database system was shut down at 2023-01-24 14:22:16 UTC
postgres_1  | 2023-01-24 14:22:30.728 UTC [1] LOG:  database system is ready to accept connections
^CGracefully stopping... (press Ctrl+C again to force)
Stopping go-webapp_postgres_1 ... done
```

```bash
 ~/tutorial/go-webapp | master ?3  docker-compose up                  ok | 16s | 23:22:45 
Creating network "go-webapp_default" with the default driver
Creating go-webapp_postgres_1 ... done
Attaching to go-webapp_postgres_1
postgres_1  | The files belonging to this database system will be owned by user "postgres".
postgres_1  | This user must also own the server process.
postgres_1  | 
postgres_1  | The database cluster will be initialized with locale "en_US.utf8".
postgres_1  | The default database encoding has accordingly been set to "UTF8".
postgres_1  | The default text search configuration will be set to "english".
postgres_1  | 
postgres_1  | Data page checksums are disabled.
postgres_1  | 
postgres_1  | fixing permissions on existing directory /var/lib/postgresql/data ... ok
postgres_1  | creating subdirectories ... ok
postgres_1  | selecting dynamic shared memory implementation ... posix
postgres_1  | selecting default max_connections ... 100
postgres_1  | selecting default shared_buffers ... 128MB
postgres_1  | selecting default time zone ... Etc/UTC
postgres_1  | creating configuration files ... ok
postgres_1  | running bootstrap script ... ok
postgres_1  | performing post-bootstrap initialization ... ok
postgres_1  | syncing data to disk ... initdb: warning: enabling "trust" authentication for local connections
postgres_1  | You can change this by editing pg_hba.conf or using the option -A, or
postgres_1  | --auth-local and --auth-host, the next time you run initdb.
postgres_1  | ok
postgres_1  | 
postgres_1  | 
postgres_1  | Success. You can now start the database server using:
postgres_1  | 
postgres_1  |     pg_ctl -D /var/lib/postgresql/data -l logfile start
postgres_1  | 
postgres_1  | waiting for server to start....2023-01-24 14:23:22.818 UTC [50] LOG:  starting PostgreSQL 14.5 (Debian 14.5-2.pgdg110+2) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
postgres_1  | 2023-01-24 14:23:22.820 UTC [50] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
postgres_1  | 2023-01-24 14:23:22.843 UTC [51] LOG:  database system was shut down at 2023-01-24 14:23:20 UTC
postgres_1  | 2023-01-24 14:23:22.876 UTC [50] LOG:  database system is ready to accept connections
postgres_1  |  done
postgres_1  | server started
postgres_1  | CREATE DATABASE
postgres_1  | 
postgres_1  | 
postgres_1  | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/create_tables.sql
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  |  set_config 
postgres_1  | ------------
postgres_1  |  
postgres_1  | (1 row)
postgres_1  | 
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  | SET
postgres_1  | CREATE TABLE
postgres_1  | ALTER TABLE
postgres_1  | CREATE TABLE
postgres_1  | ALTER TABLE
postgres_1  | COPY 0
postgres_1  | COPY 1
postgres_1  |  setval 
postgres_1  | --------
postgres_1  |       1
postgres_1  | (1 row)
postgres_1  | 
postgres_1  |  setval 
postgres_1  | --------
postgres_1  |       1
postgres_1  | (1 row)
postgres_1  | 
postgres_1  | ALTER TABLE
postgres_1  | ALTER TABLE
postgres_1  | ALTER TABLE
postgres_1  | 
postgres_1  | 
postgres_1  | 2023-01-24 14:23:26.060 UTC [50] LOG:  received fast shutdown request
postgres_1  | waiting for server to shut down...2023-01-24 14:23:26.063 UTC [50] LOG:  aborting any active transactions
postgres_1  | .2023-01-24 14:23:26.066 UTC [50] LOG:  background worker "logical replication launcher" (PID 57) exited with exit code 1
postgres_1  | 2023-01-24 14:23:26.066 UTC [52] LOG:  shutting down
postgres_1  | 2023-01-24 14:23:26.248 UTC [50] LOG:  database system is shut down
postgres_1  |  done
postgres_1  | server stopped
postgres_1  | 
postgres_1  | PostgreSQL init process complete; ready for start up.
postgres_1  | 
postgres_1  | 2023-01-24 14:23:26.308 UTC [1] LOG:  starting PostgreSQL 14.5 (Debian 14.5-2.pgdg110+2) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
postgres_1  | 2023-01-24 14:23:26.308 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
postgres_1  | 2023-01-24 14:23:26.308 UTC [1] LOG:  listening on IPv6 address "::", port 5432
postgres_1  | 2023-01-24 14:23:26.315 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
postgres_1  | 2023-01-24 14:23:26.350 UTC [66] LOG:  database system was shut down at 2023-01-24 14:23:26 UTC
postgres_1  | 2023-01-24 14:23:26.386 UTC [1] LOG:  database system is ready to accept connections
```