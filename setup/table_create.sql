
-- structure for table cushion.customer
CREATE TABLE `customer` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT COMMENT 'table row',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'customer name',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created timestamp',
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='List of customers';

-- data for table cushion.customer: ~1 rows (approximately)
INSERT INTO `customer` (`name`) VALUES
	('Josh');


-- structure for table cushion.customer_investments
CREATE TABLE `customer_investments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'table row id',
  `investment_type_id` tinyint(3) unsigned NOT NULL,
  `customer_id` mediumint(8) unsigned NOT NULL,
  `amount` mediumint(8) unsigned NOT NULL,
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- data for table cushion.customer_investments: ~0 rows (approximately)

-- structure for table cushion.investment_types
CREATE TABLE `investment_types` (
  `id` tinyint(3) unsigned NOT NULL AUTO_INCREMENT COMMENT 'row id',
  `name` varchar(255) NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='List of customers';

-- data for table cushion.investment_types: ~2 rows (approximately)
INSERT INTO `investment_types` (`name`) VALUES
	('Cushon Equities Fund'),
	('Cushon Fixed income Fund');
