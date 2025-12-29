-- Script to add hardcopy_admin_fee column to property_service_development.products table
-- Run this script if the column doesn't exist

USE property_service_development;

-- Check if hardcopy_admin_fee column exists, if not, add it
SET @dbname = DATABASE();
SET @tablename = 'products';
SET @columnname = 'hardcopy_admin_fee';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (COLUMN_NAME = @columnname)
  ) > 0,
  'SELECT "Column hardcopy_admin_fee already exists" AS result;',
  CONCAT('ALTER TABLE ', @tablename, ' ADD COLUMN ', @columnname, ' INT DEFAULT 0 AFTER coverage_type;')
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Alternative: If hardcopy_admin_fee doesn't exist, try admin_fee
SET @columnname2 = 'admin_fee';
SET @preparedStatement2 = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (COLUMN_NAME = @columnname2)
  ) > 0,
  'SELECT "Column admin_fee already exists" AS result;',
  CONCAT('ALTER TABLE ', @tablename, ' ADD COLUMN ', @columnname2, ' INT DEFAULT 0 AFTER coverage_type;')
));
PREPARE alterIfNotExists2 FROM @preparedStatement2;
EXECUTE alterIfNotExists2;
DEALLOCATE PREPARE alterIfNotExists2;

SELECT "Script completed. Check if columns were added." AS result;



