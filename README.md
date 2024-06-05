# go-mysql-diff

Compare two mysql database tables, columns and indexing tools

### 配置文件利用[toml](https://github.com/toml-lang/toml)

格式如下：

```

[servers]
  [servers.1]
  host = "127.0.0.1"
  port = "3306"
  user = "ch"
  password = "123456"
  name = "ch1"

  [servers.2]
  host = "127.0.0.1"
  port = "3306"
  user = "ch"
  password = "123456"
  name = "ch2"
```

- run

make run

- output

cat diff.sql

```sql

ALTER TABLE `xxx` ADD COLUMN `aa_id` varchar(64) NOT NULL DEFAULT '',
 ADD COLUMN `bb` decimal(32,8) NOT NULL DEFAULT '0.00000000';

ALTER TABLE `yyy` ADD COLUMN `cc` tinyint(1) NULL DEFAULT '0',
 ADD COLUMN `dd` tinyint(3) unsigned NULL DEFAULT '0',
 ADD COLUMN `aa_id` varchar(64) NOT NULL DEFAULT '';

```
