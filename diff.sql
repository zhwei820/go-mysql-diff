--- diffSql
 CREATE TABLE `fee_model` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  `name` varchar(64)  NOT NULL DEFAULT '' COMMENT 'name',
  `currency` varchar(64)  NOT NULL DEFAULT '' COMMENT 'currency',
  `is_template` tinyint(1) NOT NULL DEFAULT '0' COMMENT '''IsTemplate''',
  `type` varchar(20)  NOT NULL DEFAULT '' COMMENT 'type',
  `fee_model_id` varchar(64)  NOT NULL DEFAULT '' COMMENT 'FeeModelID',
  `mp_performance_fee_mode` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'mp performance fee mode',
  `fm_performance_fee_mode` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'fm performance fee mode',
  `mp_performance_fee_rate` text  NOT NULL COMMENT 'performance fee by mp',
  `fm_performance_fee_rate` text  NOT NULL COMMENT 'performance fee by fm',
  `mp_platform_fixed_fee_rate` text  NOT NULL COMMENT 'mp platform fixed fee rate',
  `fm_platform_fixed_fee_rate` text  NOT NULL COMMENT 'fm platform fixed fee rate',
  `redemption_fee_rate` text  NOT NULL COMMENT 'redemption fee rate',
  `subscription_fee_rate` decimal(38,18) NOT NULL DEFAULT '0.000000000000000000' COMMENT 'subscription fee rate',
  `profit_bill_cycling_month` int(11) NOT NULL DEFAULT '0' COMMENT 'profit bill cycling month',
  `profit_bill_cycling_day` int(11) NOT NULL DEFAULT '0' COMMENT 'profit bill cycling day',
  `network_group_commission_rebate` decimal(38,18) NOT NULL DEFAULT '0.000000000000000000' COMMENT 'network group commission rebate',
  `fm_profit_acc_commission_rebate` decimal(38,18) NOT NULL DEFAULT '0.000000000000000000' COMMENT 'fm profit acc commission rebate',
  `commission_rebate_scope` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'commission rebate scope',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_fee_model_id` (`fee_model_id`)
)    comment 'no qa'; 

CREATE TABLE `product_activity` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  `product_id` varchar(64)  NOT NULL DEFAULT '' COMMENT 'product id',
  `start_time` bigint(20) NOT NULL DEFAULT '0' COMMENT 'start time',
  `end_time` bigint(20) NOT NULL DEFAULT '0' COMMENT 'end time',
  `icon_url` text  COMMENT 'icon url',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT 'status',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_product_id` (`product_id`)
)    comment 'no qa'; 

CREATE TABLE `recommendation_template_product_binding` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `prioritization` bigint(20) NOT NULL DEFAULT '0' COMMENT 'prioritization',
  `product_id` varchar(64)  NOT NULL DEFAULT '' COMMENT 'product id',
  `template_id` varchar(64)  NOT NULL DEFAULT '' COMMENT 'template id',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`)
)    comment 'no qa'; 

ALTER TABLE `investment_token_config` ADD COLUMN `standard_fee_model_id` varchar(64) NOT NULL DEFAULT '',
 ADD COLUMN `advanced_fee_model_id` varchar(64) NOT NULL DEFAULT '';

ALTER TABLE `private_portfolio` ADD COLUMN `fee_model_id` varchar(64) NOT NULL DEFAULT '',
 ADD COLUMN `expected_apy` decimal(32,8) NOT NULL DEFAULT '0.00000000';

ALTER TABLE `product` ADD COLUMN `fee_model_id` varchar(64) NOT NULL DEFAULT '',
 ADD COLUMN `default_shown_nav` tinyint(3) unsigned NULL DEFAULT '0',
 ADD COLUMN `ignore_data_cache_update` tinyint(1) NULL DEFAULT '0';

ALTER TABLE `product_content` ADD COLUMN `currency_positions_analysis_for_investor` int(10) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `trade_history_positions_analysis_for_code_holder` int(10) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `performance_display_pnl` tinyint(3) unsigned NOT NULL DEFAULT '2',
 ADD COLUMN `private_strategy_visibility` bigint(20) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `currency_positions_analysis_for_code_holder` int(10) unsigned NOT NULL DEFAULT '0',
 ADD COLUMN `trade_history_positions_analysis_for_investor` int(10) unsigned NOT NULL DEFAULT '0';

ALTER TABLE `access_code` ADD INDEX `idx_portfolio_id` (portfolio_id,product_id),
 ADD UNIQUE INDEX `idx_user_id` (user_id,code);

ALTER TABLE `product_content` ADD INDEX `idx_product_content_deleted_at` (deleted_at);


