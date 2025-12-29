-- Script to alter hardcopy_admin_fee column from DECIMAL(10,2) to INT
-- Run this script to change existing column type

USE property_service_development;

-- Check if hardcopy_admin_fee column exists and alter it to INT
SET @dbname = DATABASE();
SET @tablename = 'products';
SET @columnname = 'hardcopy_admin_fee';

-- First, check if column exists
SET @columnExists = (
  SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE
    (TABLE_SCHEMA = @dbname)
    AND (TABLE_NAME = @tablename)
    AND (COLUMN_NAME = @columnname)
);

-- If column exists, alter it to INT
SET @preparedStatement = (SELECT IF(
  @columnExists > 0,
  CONCAT('ALTER TABLE ', @tablename, ' MODIFY COLUMN ', @columnname, ' INT DEFAULT 0;'),
  'SELECT "Column hardcopy_admin_fee does not exist" AS result;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

-- Also check and alter admin_fee column if it exists
SET @columnname2 = 'admin_fee';
SET @columnExists2 = (
  SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE
    (TABLE_SCHEMA = @dbname)
    AND (TABLE_NAME = @tablename)
    AND (COLUMN_NAME = @columnname2)
);

SET @preparedStatement2 = (SELECT IF(
  @columnExists2 > 0,
  CONCAT('ALTER TABLE ', @tablename, ' MODIFY COLUMN ', @columnname2, ' INT DEFAULT 0;'),
  'SELECT "Column admin_fee does not exist" AS result;'
));
PREPARE alterIfExists2 FROM @preparedStatement2;
EXECUTE alterIfExists2;
DEALLOCATE PREPARE alterIfExists2;

SELECT "Script completed. Check if columns were altered." AS result;


