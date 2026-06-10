-- Rollback: 013_create_incentive_tables

DROP INDEX IF EXISTS idx_staff_incentives_status;
DROP INDEX IF EXISTS idx_staff_incentives_period;
DROP INDEX IF EXISTS idx_staff_incentives_rule_id;
DROP INDEX IF EXISTS idx_staff_incentives_staff_id;
DROP TABLE IF EXISTS staff_incentives;

DROP INDEX IF EXISTS idx_incentive_rule_slabs_rule_id;
DROP TABLE IF EXISTS incentive_rule_slabs;

DROP INDEX IF EXISTS idx_incentive_rules_deleted_at;
DROP INDEX IF EXISTS idx_incentive_rules_staff_id;
DROP INDEX IF EXISTS idx_incentive_rules_is_active;
DROP INDEX IF EXISTS idx_incentive_rules_type;
DROP TABLE IF EXISTS incentive_rules;
