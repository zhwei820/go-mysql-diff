# go-mysql-diff

Tools for Comparing two mysql database tables, columns and indexes. Only support add table/column/index.

### Edit Config file

```
vi config.toml
```

```toml
# Only supports adding tables/fields/indexes
[servers.1] # from
host = "127.0.0.1"
port = "3316"
user = "root"
password = "123456"
name = "ch_main"

[servers.2] # to
host = "127.0.0.1"
port = "3316"
user = "root"
password = "123456"
name = "ch_v2"

```

### Run to make a diff

```bash
make run
```

### Check output

```bash
> cat diff.sql
```

```sql

ALTER TABLE `xxx` ADD COLUMN `aa_id` varchar(64) NOT NULL DEFAULT '',
 ADD COLUMN `bb` decimal(32,8) NOT NULL DEFAULT '0.00000000';

ALTER TABLE `yyy` ADD COLUMN `cc` tinyint(1) NULL DEFAULT '0',
 ADD COLUMN `dd` tinyint(3) unsigned NULL DEFAULT '0',
 ADD COLUMN `aa_id` varchar(64) NOT NULL DEFAULT '';

ALTER TABLE `zzz` ADD INDEX `idx_iii` (deleted_at);

```
