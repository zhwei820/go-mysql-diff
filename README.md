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

ALTER TABLE `private_portfolio` ADD COLUMN `fee_model_id` varchar(64) NOT NULL DEFAULT '',
 ADD COLUMN `expected_apy` decimal(32,8) NOT NULL DEFAULT '0.00000000';

ALTER TABLE `product` ADD COLUMN `ignore_data_cache_update` tinyint(1) NULL DEFAULT '0',
 ADD COLUMN `default_shown_nav` tinyint(3) unsigned NULL DEFAULT '0',
 ADD COLUMN `fee_model_id` varchar(64) NOT NULL DEFAULT '';

ALTER TABLE `product_content` ADD COLUMN `currency_positions_analysis_for_code_holder` int(10) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `trade_history_positions_analysis_for_investor` int(10) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `performance_display_pnl` tinyint(3) unsigned NOT NULL DEFAULT '2',
 ADD COLUMN `private_strategy_visibility` bigint(20) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `trade_history_positions_analysis_for_code_holder` int(10) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `currency_positions_analysis_for_investor` int(10) unsigned NOT NULL DEFAULT '0';

```
