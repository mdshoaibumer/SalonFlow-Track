export namespace domain {
	
	export class Advance {
	    id: number[];
	    staff_id: number[];
	    staff_name?: string;
	    amount: number;
	    advance_date: string;
	    reason: string;
	    recovered_amount: number;
	    remaining_amount: number;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Advance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.staff_id = source["staff_id"];
	        this.staff_name = source["staff_name"];
	        this.amount = source["amount"];
	        this.advance_date = source["advance_date"];
	        this.reason = source["reason"];
	        this.recovered_amount = source["recovered_amount"];
	        this.remaining_amount = source["remaining_amount"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppointmentService {
	    id: number[];
	    appointment_id: number[];
	    service_id: string;
	    service_name: string;
	    duration_minutes: number;
	    price: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new AppointmentService(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.appointment_id = source["appointment_id"];
	        this.service_id = source["service_id"];
	        this.service_name = source["service_name"];
	        this.duration_minutes = source["duration_minutes"];
	        this.price = source["price"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Appointment {
	    id: number[];
	    customer_id: string;
	    staff_id: string;
	    appointment_date: string;
	    start_time: string;
	    end_time: string;
	    status: string;
	    notes: string;
	    is_walkin: boolean;
	    total_amount: number;
	    services?: AppointmentService[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Appointment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.customer_id = source["customer_id"];
	        this.staff_id = source["staff_id"];
	        this.appointment_date = source["appointment_date"];
	        this.start_time = source["start_time"];
	        this.end_time = source["end_time"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.is_walkin = source["is_walkin"];
	        this.total_amount = source["total_amount"];
	        this.services = this.convertValues(source["services"], AppointmentService);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppointmentFilter {
	    date: string;
	    staff_id: string;
	    customer_id: string;
	    status: string;
	    start_date: string;
	    end_date: string;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new AppointmentFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.staff_id = source["staff_id"];
	        this.customer_id = source["customer_id"];
	        this.status = source["status"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class AppointmentHistory {
	    id: number[];
	    appointment_id: number[];
	    old_status: string;
	    new_status: string;
	    changed_by: string;
	    note: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new AppointmentHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.appointment_id = source["appointment_id"];
	        this.old_status = source["old_status"];
	        this.new_status = source["new_status"];
	        this.changed_by = source["changed_by"];
	        this.note = source["note"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class AutomationRule {
	    id: number[];
	    name: string;
	    trigger_type: string;
	    template_id: string;
	    delay_minutes: number;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new AutomationRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.trigger_type = source["trigger_type"];
	        this.template_id = source["template_id"];
	        this.delay_minutes = source["delay_minutes"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BackupRecord {
	    id: number[];
	    backup_name: string;
	    backup_type: string;
	    backup_path: string;
	    file_size: number;
	    checksum: string;
	    status: string;
	    error_message?: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new BackupRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.backup_name = source["backup_name"];
	        this.backup_type = source["backup_type"];
	        this.backup_path = source["backup_path"];
	        this.file_size = source["file_size"];
	        this.checksum = source["checksum"];
	        this.status = source["status"];
	        this.error_message = source["error_message"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BackupStats {
	    total_backups: number;
	    last_backup_name: string;
	    last_backup_date: string;
	    last_backup_size: number;
	    last_status: string;
	    total_restores: number;
	
	    static createFrom(source: any = {}) {
	        return new BackupStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_backups = source["total_backups"];
	        this.last_backup_name = source["last_backup_name"];
	        this.last_backup_date = source["last_backup_date"];
	        this.last_backup_size = source["last_backup_size"];
	        this.last_status = source["last_status"];
	        this.total_restores = source["total_restores"];
	    }
	}
	export class BackupVerification {
	    backup_id: string;
	    file_exists: boolean;
	    can_open: boolean;
	    integrity_ok: boolean;
	    checksum_ok: boolean;
	    status: string;
	    error_message?: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupVerification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.backup_id = source["backup_id"];
	        this.file_exists = source["file_exists"];
	        this.can_open = source["can_open"];
	        this.integrity_ok = source["integrity_ok"];
	        this.checksum_ok = source["checksum_ok"];
	        this.status = source["status"];
	        this.error_message = source["error_message"];
	    }
	}
	export class CategoryExpense {
	    category_id: string;
	    category_name: string;
	    amount: number;
	    percentage: number;
	
	    static createFrom(source: any = {}) {
	        return new CategoryExpense(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category_id = source["category_id"];
	        this.category_name = source["category_name"];
	        this.amount = source["amount"];
	        this.percentage = source["percentage"];
	    }
	}
	export class CloudBackupConfig {
	    id: number[];
	    provider: string;
	    bucket_name: string;
	    region: string;
	    access_key: string;
	    endpoint: string;
	    encrypt_backups: boolean;
	    auto_backup: boolean;
	    auto_backup_interval_hours: number;
	    max_versions: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new CloudBackupConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.provider = source["provider"];
	        this.bucket_name = source["bucket_name"];
	        this.region = source["region"];
	        this.access_key = source["access_key"];
	        this.endpoint = source["endpoint"];
	        this.encrypt_backups = source["encrypt_backups"];
	        this.auto_backup = source["auto_backup"];
	        this.auto_backup_interval_hours = source["auto_backup_interval_hours"];
	        this.max_versions = source["max_versions"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CloudBackupHistory {
	    id: number[];
	    provider: string;
	    file_name: string;
	    file_size: number;
	    remote_path: string;
	    status: string;
	    is_encrypted: boolean;
	    error_message: string;
	    started_at: string;
	    completed_at: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new CloudBackupHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.provider = source["provider"];
	        this.file_name = source["file_name"];
	        this.file_size = source["file_size"];
	        this.remote_path = source["remote_path"];
	        this.status = source["status"];
	        this.is_encrypted = source["is_encrypted"];
	        this.error_message = source["error_message"];
	        this.started_at = source["started_at"];
	        this.completed_at = source["completed_at"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CloudBackupStats {
	    last_backup_at: string;
	    total_backups: number;
	    total_size_bytes: number;
	    provider: string;
	    auto_enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CloudBackupStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.last_backup_at = source["last_backup_at"];
	        this.total_backups = source["total_backups"];
	        this.total_size_bytes = source["total_size_bytes"];
	        this.provider = source["provider"];
	        this.auto_enabled = source["auto_enabled"];
	    }
	}
	export class ColumnMapping {
	    source_column: string;
	    target_field: string;
	
	    static createFrom(source: any = {}) {
	        return new ColumnMapping(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_column = source["source_column"];
	        this.target_field = source["target_field"];
	    }
	}
	export class CommissionRule {
	    id: number[];
	    rule_name: string;
	    rule_type: string;
	    target_type: string;
	    target_id: string;
	    calculation_type: string;
	    calculation_value: number;
	    minimum_target: number;
	    maximum_target: number;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new CommissionRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.rule_name = source["rule_name"];
	        this.rule_type = source["rule_type"];
	        this.target_type = source["target_type"];
	        this.target_id = source["target_id"];
	        this.calculation_type = source["calculation_type"];
	        this.calculation_value = source["calculation_value"];
	        this.minimum_target = source["minimum_target"];
	        this.maximum_target = source["maximum_target"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CommissionTransaction {
	    id: number[];
	    staff_id: number[];
	    invoice_id: number[];
	    rule_id: number[];
	    revenue_amount: number;
	    commission_amount: number;
	    business_date: string;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new CommissionTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.staff_id = source["staff_id"];
	        this.invoice_id = source["invoice_id"];
	        this.rule_id = source["rule_id"];
	        this.revenue_amount = source["revenue_amount"];
	        this.commission_amount = source["commission_amount"];
	        this.business_date = source["business_date"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Customer {
	    id: number[];
	    customer_code: string;
	    full_name: string;
	    phone: string;
	    email: string;
	    gender: string;
	    // Go type: time
	    date_of_birth?: any;
	    // Go type: time
	    anniversary_date?: any;
	    address: string;
	    notes: string;
	    total_visits: number;
	    total_spent: number;
	    // Go type: time
	    last_visit_date?: any;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Customer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.customer_code = source["customer_code"];
	        this.full_name = source["full_name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gender = source["gender"];
	        this.date_of_birth = this.convertValues(source["date_of_birth"], null);
	        this.anniversary_date = this.convertValues(source["anniversary_date"], null);
	        this.address = source["address"];
	        this.notes = source["notes"];
	        this.total_visits = source["total_visits"];
	        this.total_spent = source["total_spent"];
	        this.last_visit_date = this.convertValues(source["last_visit_date"], null);
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrendPoint {
	    period: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new TrendPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.value = source["value"];
	    }
	}
	export class NameValuePair {
	    name: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new NameValuePair(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class CustomerReport {
	    total_customers: number;
	    new_customers: number;
	    repeat_customers: number;
	    birthday_today: number;
	    inactive_count: number;
	    top_customers: NameValuePair[];
	    growth_trend: TrendPoint[];
	
	    static createFrom(source: any = {}) {
	        return new CustomerReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_customers = source["total_customers"];
	        this.new_customers = source["new_customers"];
	        this.repeat_customers = source["repeat_customers"];
	        this.birthday_today = source["birthday_today"];
	        this.inactive_count = source["inactive_count"];
	        this.top_customers = this.convertValues(source["top_customers"], NameValuePair);
	        this.growth_trend = this.convertValues(source["growth_trend"], TrendPoint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DashboardStats {
	    today_revenue: number;
	    today_customers: number;
	    today_invoices: number;
	    monthly_revenue: number;
	    monthly_expenses: number;
	    monthly_profit: number;
	    inventory_value: number;
	    outstanding_salary: number;
	    outstanding_advances: number;
	    low_stock_count: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.today_revenue = source["today_revenue"];
	        this.today_customers = source["today_customers"];
	        this.today_invoices = source["today_invoices"];
	        this.monthly_revenue = source["monthly_revenue"];
	        this.monthly_expenses = source["monthly_expenses"];
	        this.monthly_profit = source["monthly_profit"];
	        this.inventory_value = source["inventory_value"];
	        this.outstanding_salary = source["outstanding_salary"];
	        this.outstanding_advances = source["outstanding_advances"];
	        this.low_stock_count = source["low_stock_count"];
	    }
	}
	export class DualTrendPoint {
	    period: string;
	    value1: number;
	    value2: number;
	
	    static createFrom(source: any = {}) {
	        return new DualTrendPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.value1 = source["value1"];
	        this.value2 = source["value2"];
	    }
	}
	export class Expense {
	    id: number[];
	    expense_number: string;
	    category_id: number[];
	    category_name?: string;
	    amount: number;
	    expense_date: string;
	    payment_method: string;
	    vendor_name: string;
	    invoice_reference: string;
	    description: string;
	    attachment_path: string;
	    status: string;
	    created_by: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Expense(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.expense_number = source["expense_number"];
	        this.category_id = source["category_id"];
	        this.category_name = source["category_name"];
	        this.amount = source["amount"];
	        this.expense_date = source["expense_date"];
	        this.payment_method = source["payment_method"];
	        this.vendor_name = source["vendor_name"];
	        this.invoice_reference = source["invoice_reference"];
	        this.description = source["description"];
	        this.attachment_path = source["attachment_path"];
	        this.status = source["status"];
	        this.created_by = source["created_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExpenseCategory {
	    id: number[];
	    name: string;
	    description: string;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExpenseReport {
	    total_expenses: number;
	    by_category: NameValuePair[];
	    monthly_trend: TrendPoint[];
	    revenue_vs_expense: DualTrendPoint[];
	
	    static createFrom(source: any = {}) {
	        return new ExpenseReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_expenses = source["total_expenses"];
	        this.by_category = this.convertValues(source["by_category"], NameValuePair);
	        this.monthly_trend = this.convertValues(source["monthly_trend"], TrendPoint);
	        this.revenue_vs_expense = this.convertValues(source["revenue_vs_expense"], DualTrendPoint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExpenseStats {
	    today_expenses: number;
	    monthly_expenses: number;
	    today_revenue: number;
	    monthly_revenue: number;
	    monthly_profit: number;
	    profit_margin: number;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.today_expenses = source["today_expenses"];
	        this.monthly_expenses = source["monthly_expenses"];
	        this.today_revenue = source["today_revenue"];
	        this.monthly_revenue = source["monthly_revenue"];
	        this.monthly_profit = source["monthly_profit"];
	        this.profit_margin = source["profit_margin"];
	    }
	}
	export class GSTReport {
	    period: string;
	    start_date: string;
	    end_date: string;
	    total_invoices: number;
	    taxable_amount: number;
	    total_cgst: number;
	    total_sgst: number;
	    total_igst: number;
	    total_tax: number;
	    grand_total: number;
	
	    static createFrom(source: any = {}) {
	        return new GSTReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.total_invoices = source["total_invoices"];
	        this.taxable_amount = source["taxable_amount"];
	        this.total_cgst = source["total_cgst"];
	        this.total_sgst = source["total_sgst"];
	        this.total_igst = source["total_igst"];
	        this.total_tax = source["total_tax"];
	        this.grand_total = source["grand_total"];
	    }
	}
	export class GSTReportFilter {
	    period: string;
	    start_date: string;
	    end_date: string;
	
	    static createFrom(source: any = {}) {
	        return new GSTReportFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	    }
	}
	export class GSTSettings {
	    id: number[];
	    business_name: string;
	    gstin: string;
	    state: string;
	    address: string;
	    hsn_code: string;
	    cgst_rate: number;
	    sgst_rate: number;
	    igst_rate: number;
	    is_gst_enabled: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new GSTSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.business_name = source["business_name"];
	        this.gstin = source["gstin"];
	        this.state = source["state"];
	        this.address = source["address"];
	        this.hsn_code = source["hsn_code"];
	        this.cgst_rate = source["cgst_rate"];
	        this.sgst_rate = source["sgst_rate"];
	        this.igst_rate = source["igst_rate"];
	        this.is_gst_enabled = source["is_gst_enabled"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportJob {
	    id: number[];
	    template_id?: string;
	    file_name: string;
	    file_path: string;
	    target_entity: string;
	    status: string;
	    total_rows: number;
	    valid_rows: number;
	    invalid_rows: number;
	    imported_rows: number;
	    column_mapping: string;
	    error_message?: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ImportJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.template_id = source["template_id"];
	        this.file_name = source["file_name"];
	        this.file_path = source["file_path"];
	        this.target_entity = source["target_entity"];
	        this.status = source["status"];
	        this.total_rows = source["total_rows"];
	        this.valid_rows = source["valid_rows"];
	        this.invalid_rows = source["invalid_rows"];
	        this.imported_rows = source["imported_rows"];
	        this.column_mapping = source["column_mapping"];
	        this.error_message = source["error_message"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportLog {
	    id: number[];
	    job_id: number[];
	    row_number: number;
	    status: string;
	    message: string;
	    row_data: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ImportLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.job_id = source["job_id"];
	        this.row_number = source["row_number"];
	        this.status = source["status"];
	        this.message = source["message"];
	        this.row_data = source["row_data"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportLogRow {
	    row_number: number;
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportLogRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.row_number = source["row_number"];
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}
	export class ImportPreview {
	    job_id: number[];
	    total_rows: number;
	    valid_rows: number;
	    invalid_rows: number;
	    warnings: number;
	    headers: string[];
	    sample_rows: string[][];
	    errors?: ImportLogRow[];
	
	    static createFrom(source: any = {}) {
	        return new ImportPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.job_id = source["job_id"];
	        this.total_rows = source["total_rows"];
	        this.valid_rows = source["valid_rows"];
	        this.invalid_rows = source["invalid_rows"];
	        this.warnings = source["warnings"];
	        this.headers = source["headers"];
	        this.sample_rows = source["sample_rows"];
	        this.errors = this.convertValues(source["errors"], ImportLogRow);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InventoryReport {
	    total_value: number;
	    low_stock_count: number;
	    fast_moving: NameValuePair[];
	    slow_moving: NameValuePair[];
	    purchase_trend: TrendPoint[];
	    consumption_trend: TrendPoint[];
	
	    static createFrom(source: any = {}) {
	        return new InventoryReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_value = source["total_value"];
	        this.low_stock_count = source["low_stock_count"];
	        this.fast_moving = this.convertValues(source["fast_moving"], NameValuePair);
	        this.slow_moving = this.convertValues(source["slow_moving"], NameValuePair);
	        this.purchase_trend = this.convertValues(source["purchase_trend"], TrendPoint);
	        this.consumption_trend = this.convertValues(source["consumption_trend"], TrendPoint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InventoryStats {
	    total_products: number;
	    active_products: number;
	    low_stock_count: number;
	    total_value: number;
	    total_purchases_this_month: number;
	
	    static createFrom(source: any = {}) {
	        return new InventoryStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_products = source["total_products"];
	        this.active_products = source["active_products"];
	        this.low_stock_count = source["low_stock_count"];
	        this.total_value = source["total_value"];
	        this.total_purchases_this_month = source["total_purchases_this_month"];
	    }
	}
	export class InvoiceItem {
	    id: number[];
	    invoice_id: number[];
	    service_id: number[];
	    service_name_snapshot: string;
	    quantity: number;
	    unit_price: number;
	    discount: number;
	    line_total: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_id = source["invoice_id"];
	        this.service_id = source["service_id"];
	        this.service_name_snapshot = source["service_name_snapshot"];
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	        this.discount = source["discount"];
	        this.line_total = source["line_total"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Invoice {
	    id: number[];
	    invoice_number: string;
	    customer_id: number[];
	    staff_id: number[];
	    items: InvoiceItem[];
	    subtotal: number;
	    discount: number;
	    tax: number;
	    grand_total: number;
	    payment_status: string;
	    payment_method: string;
	    notes: string;
	    // Go type: time
	    invoice_date: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Invoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_number = source["invoice_number"];
	        this.customer_id = source["customer_id"];
	        this.staff_id = source["staff_id"];
	        this.items = this.convertValues(source["items"], InvoiceItem);
	        this.subtotal = source["subtotal"];
	        this.discount = source["discount"];
	        this.tax = source["tax"];
	        this.grand_total = source["grand_total"];
	        this.payment_status = source["payment_status"];
	        this.payment_method = source["payment_method"];
	        this.notes = source["notes"];
	        this.invoice_date = this.convertValues(source["invoice_date"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class KPIMetrics {
	    revenue_growth_pct: number;
	    customer_growth_pct: number;
	    profit_margin_pct: number;
	    average_bill_value: number;
	    repeat_customer_pct: number;
	    staff_productivity_pct: number;
	
	    static createFrom(source: any = {}) {
	        return new KPIMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.revenue_growth_pct = source["revenue_growth_pct"];
	        this.customer_growth_pct = source["customer_growth_pct"];
	        this.profit_margin_pct = source["profit_margin_pct"];
	        this.average_bill_value = source["average_bill_value"];
	        this.repeat_customer_pct = source["repeat_customer_pct"];
	        this.staff_productivity_pct = source["staff_productivity_pct"];
	    }
	}
	export class License {
	    id: number[];
	    license_key: string;
	    customer_name: string;
	    salon_name: string;
	    device_id: string;
	    issued_date: string;
	    expiry_date: string;
	    status: string;
	    signature: string;
	    last_validation: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new License(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.license_key = source["license_key"];
	        this.customer_name = source["customer_name"];
	        this.salon_name = source["salon_name"];
	        this.device_id = source["device_id"];
	        this.issued_date = source["issued_date"];
	        this.expiry_date = source["expiry_date"];
	        this.status = source["status"];
	        this.signature = source["signature"];
	        this.last_validation = source["last_validation"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LicenseEvent {
	    id: number[];
	    license_id: number[];
	    event_type: string;
	    event_date: string;
	    notes: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new LicenseEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.license_id = source["license_id"];
	        this.event_type = source["event_type"];
	        this.event_date = source["event_date"];
	        this.notes = source["notes"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LicenseStatus {
	    license?: License;
	    days_remaining: number;
	    grace_days_remaining: number;
	    is_restricted: boolean;
	    needs_renewal: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LicenseStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.license = this.convertValues(source["license"], License);
	        this.days_remaining = source["days_remaining"];
	        this.grace_days_remaining = source["grace_days_remaining"];
	        this.is_restricted = source["is_restricted"];
	        this.needs_renewal = source["needs_renewal"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LicenseValidation {
	    valid: boolean;
	    status: string;
	    days_remaining: number;
	    is_restricted: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LicenseValidation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.status = source["status"];
	        this.days_remaining = source["days_remaining"];
	        this.is_restricted = source["is_restricted"];
	        this.message = source["message"];
	    }
	}
	export class LowStockItem {
	    product_id: string;
	    product_code: string;
	    product_name: string;
	    category: string;
	    current_stock: number;
	    minimum_stock: number;
	    deficit: number;
	
	    static createFrom(source: any = {}) {
	        return new LowStockItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.product_code = source["product_code"];
	        this.product_name = source["product_name"];
	        this.category = source["category"];
	        this.current_stock = source["current_stock"];
	        this.minimum_stock = source["minimum_stock"];
	        this.deficit = source["deficit"];
	    }
	}
	export class MemberSubscription {
	    id: number[];
	    customer_id: string;
	    plan_id: number[];
	    plan_name: string;
	    start_date: string;
	    end_date: string;
	    total_sessions: number;
	    used_sessions: number;
	    amount_paid: number;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new MemberSubscription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.customer_id = source["customer_id"];
	        this.plan_id = source["plan_id"];
	        this.plan_name = source["plan_name"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.total_sessions = source["total_sessions"];
	        this.used_sessions = source["used_sessions"];
	        this.amount_paid = source["amount_paid"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PackageService {
	    id: number[];
	    plan_id: number[];
	    service_id: string;
	    service_name: string;
	    sessions_included: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PackageService(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.plan_id = source["plan_id"];
	        this.service_id = source["service_id"];
	        this.service_name = source["service_name"];
	        this.sessions_included = source["sessions_included"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MembershipPlan {
	    id: number[];
	    name: string;
	    description: string;
	    plan_type: string;
	    price: number;
	    duration_days: number;
	    max_sessions: number;
	    discount_percentage: number;
	    priority_booking: boolean;
	    is_active: boolean;
	    services?: PackageService[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new MembershipPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.plan_type = source["plan_type"];
	        this.price = source["price"];
	        this.duration_days = source["duration_days"];
	        this.max_sessions = source["max_sessions"];
	        this.discount_percentage = source["discount_percentage"];
	        this.priority_booking = source["priority_booking"];
	        this.is_active = source["is_active"];
	        this.services = this.convertValues(source["services"], PackageService);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MembershipStats {
	    active_members: number;
	    expired_members: number;
	    total_revenue: number;
	
	    static createFrom(source: any = {}) {
	        return new MembershipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.active_members = source["active_members"];
	        this.expired_members = source["expired_members"];
	        this.total_revenue = source["total_revenue"];
	    }
	}
	export class MonthlyTrend {
	    month: string;
	    revenue: number;
	    expenses: number;
	    profit: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyTrend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.revenue = source["revenue"];
	        this.expenses = source["expenses"];
	        this.profit = source["profit"];
	    }
	}
	
	export class PLTrendPoint {
	    period: string;
	    revenue: number;
	    expenses: number;
	    salary: number;
	    profit: number;
	
	    static createFrom(source: any = {}) {
	        return new PLTrendPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.revenue = source["revenue"];
	        this.expenses = source["expenses"];
	        this.salary = source["salary"];
	        this.profit = source["profit"];
	    }
	}
	
	export class Payment {
	    id: number[];
	    invoice_id: number[];
	    amount: number;
	    payment_method: string;
	    reference_number: string;
	    // Go type: time
	    payment_date: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Payment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_id = source["invoice_id"];
	        this.amount = source["amount"];
	        this.payment_method = source["payment_method"];
	        this.reference_number = source["reference_number"];
	        this.payment_date = this.convertValues(source["payment_date"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PrintJob {
	    id: number[];
	    document_type: string;
	    document_id: string;
	    printer_name: string;
	    paper_width: string;
	    status: string;
	    copies: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PrintJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.document_type = source["document_type"];
	        this.document_id = source["document_id"];
	        this.printer_name = source["printer_name"];
	        this.paper_width = source["paper_width"];
	        this.status = source["status"];
	        this.copies = source["copies"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PrinterSettings {
	    id: number[];
	    default_printer: string;
	    paper_width: string;
	    margin_top: number;
	    margin_bottom: number;
	    margin_left: number;
	    margin_right: number;
	    header_text: string;
	    footer_text: string;
	    show_logo: boolean;
	    show_qr: boolean;
	    upi_id: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PrinterSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.default_printer = source["default_printer"];
	        this.paper_width = source["paper_width"];
	        this.margin_top = source["margin_top"];
	        this.margin_bottom = source["margin_bottom"];
	        this.margin_left = source["margin_left"];
	        this.margin_right = source["margin_right"];
	        this.header_text = source["header_text"];
	        this.footer_text = source["footer_text"];
	        this.show_logo = source["show_logo"];
	        this.show_qr = source["show_qr"];
	        this.upi_id = source["upi_id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Product {
	    id: number[];
	    product_code: string;
	    name: string;
	    category: string;
	    brand: string;
	    unit: string;
	    sku: string;
	    purchase_price: number;
	    selling_price: number;
	    current_stock: number;
	    minimum_stock: number;
	    maximum_stock: number;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_code = source["product_code"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.brand = source["brand"];
	        this.unit = source["unit"];
	        this.sku = source["sku"];
	        this.purchase_price = source["purchase_price"];
	        this.selling_price = source["selling_price"];
	        this.current_stock = source["current_stock"];
	        this.minimum_stock = source["minimum_stock"];
	        this.maximum_stock = source["maximum_stock"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProfitLoss {
	    period: string;
	    total_revenue: number;
	    total_expenses: number;
	    gross_profit: number;
	    profit_margin: number;
	    expenses_by_category: CategoryExpense[];
	
	    static createFrom(source: any = {}) {
	        return new ProfitLoss(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.total_revenue = source["total_revenue"];
	        this.total_expenses = source["total_expenses"];
	        this.gross_profit = source["gross_profit"];
	        this.profit_margin = source["profit_margin"];
	        this.expenses_by_category = this.convertValues(source["expenses_by_category"], CategoryExpense);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProfitLossReport {
	    revenue: number;
	    expenses: number;
	    salary_cost: number;
	    net_profit: number;
	    trend: PLTrendPoint[];
	
	    static createFrom(source: any = {}) {
	        return new ProfitLossReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.revenue = source["revenue"];
	        this.expenses = source["expenses"];
	        this.salary_cost = source["salary_cost"];
	        this.net_profit = source["net_profit"];
	        this.trend = this.convertValues(source["trend"], PLTrendPoint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PurchaseItem {
	    id: number[];
	    purchase_entry_id: number[];
	    product_id: number[];
	    product_name?: string;
	    quantity: number;
	    unit_price: number;
	    line_total: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.purchase_entry_id = source["purchase_entry_id"];
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	        this.line_total = source["line_total"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PurchaseEntry {
	    id: number[];
	    purchase_number: string;
	    vendor_name: string;
	    invoice_number: string;
	    purchase_date: string;
	    total_amount: number;
	    notes: string;
	    items?: PurchaseItem[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.purchase_number = source["purchase_number"];
	        this.vendor_name = source["vendor_name"];
	        this.invoice_number = source["invoice_number"];
	        this.purchase_date = source["purchase_date"];
	        this.total_amount = source["total_amount"];
	        this.notes = source["notes"];
	        this.items = this.convertValues(source["items"], PurchaseItem);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ReceiptItem {
	    name: string;
	    quantity: number;
	    price: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ReceiptItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.quantity = source["quantity"];
	        this.price = source["price"];
	        this.total = source["total"];
	    }
	}
	export class ReceiptData {
	    salon_name: string;
	    gstin: string;
	    address: string;
	    invoice_number: string;
	    date: string;
	    customer_name: string;
	    customer_phone: string;
	    items: ReceiptItem[];
	    subtotal: number;
	    cgst: number;
	    sgst: number;
	    igst: number;
	    discount: number;
	    grand_total: number;
	    payment_method: string;
	    footer_text: string;
	    upi_id: string;
	
	    static createFrom(source: any = {}) {
	        return new ReceiptData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.salon_name = source["salon_name"];
	        this.gstin = source["gstin"];
	        this.address = source["address"];
	        this.invoice_number = source["invoice_number"];
	        this.date = source["date"];
	        this.customer_name = source["customer_name"];
	        this.customer_phone = source["customer_phone"];
	        this.items = this.convertValues(source["items"], ReceiptItem);
	        this.subtotal = source["subtotal"];
	        this.cgst = source["cgst"];
	        this.sgst = source["sgst"];
	        this.igst = source["igst"];
	        this.discount = source["discount"];
	        this.grand_total = source["grand_total"];
	        this.payment_method = source["payment_method"];
	        this.footer_text = source["footer_text"];
	        this.upi_id = source["upi_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class RestoreRecord {
	    id: number[];
	    backup_id: number[];
	    backup_name: string;
	    restore_date: string;
	    status: string;
	    notes: string;
	    error_message?: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new RestoreRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.backup_id = source["backup_id"];
	        this.backup_name = source["backup_name"];
	        this.restore_date = source["restore_date"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.error_message = source["error_message"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RevenueTrendPoint {
	    date: string;
	    revenue: number;
	
	    static createFrom(source: any = {}) {
	        return new RevenueTrendPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.revenue = source["revenue"];
	    }
	}
	export class RevenueReport {
	    trend: RevenueTrendPoint[];
	    by_service: NameValuePair[];
	    by_staff: NameValuePair[];
	    by_customer: NameValuePair[];
	    total_revenue: number;
	    invoice_count: number;
	
	    static createFrom(source: any = {}) {
	        return new RevenueReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.trend = this.convertValues(source["trend"], RevenueTrendPoint);
	        this.by_service = this.convertValues(source["by_service"], NameValuePair);
	        this.by_staff = this.convertValues(source["by_staff"], NameValuePair);
	        this.by_customer = this.convertValues(source["by_customer"], NameValuePair);
	        this.total_revenue = source["total_revenue"];
	        this.invoice_count = source["invoice_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SalaryCycle {
	    id: number[];
	    month: number;
	    year: number;
	    status: string;
	    generated_at: string;
	    generated_by: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SalaryCycle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.month = source["month"];
	        this.year = source["year"];
	        this.status = source["status"];
	        this.generated_at = source["generated_at"];
	        this.generated_by = source["generated_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SalaryRecord {
	    id: number[];
	    salary_cycle_id: number[];
	    staff_id: number[];
	    staff_name?: string;
	    base_salary: number;
	    commission_amount: number;
	    bonus_amount: number;
	    advance_amount: number;
	    deduction_amount: number;
	    gross_salary: number;
	    net_salary: number;
	    payment_status: string;
	    payment_date: string;
	    notes: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SalaryRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.salary_cycle_id = source["salary_cycle_id"];
	        this.staff_id = source["staff_id"];
	        this.staff_name = source["staff_name"];
	        this.base_salary = source["base_salary"];
	        this.commission_amount = source["commission_amount"];
	        this.bonus_amount = source["bonus_amount"];
	        this.advance_amount = source["advance_amount"];
	        this.deduction_amount = source["deduction_amount"];
	        this.gross_salary = source["gross_salary"];
	        this.net_salary = source["net_salary"];
	        this.payment_status = source["payment_status"];
	        this.payment_date = source["payment_date"];
	        this.notes = source["notes"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Service {
	    id: number[];
	    service_code: string;
	    name: string;
	    category: string;
	    description: string;
	    duration_minutes: number;
	    price: number;
	    cost_price: number;
	    commission_type: string;
	    commission_value: number;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Service(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.service_code = source["service_code"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.duration_minutes = source["duration_minutes"];
	        this.price = source["price"];
	        this.cost_price = source["cost_price"];
	        this.commission_type = source["commission_type"];
	        this.commission_value = source["commission_value"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ServiceReport {
	    top_services: NameValuePair[];
	    least_used: NameValuePair[];
	    revenue_by_service: NameValuePair[];
	    avg_service_value: number;
	    total_bookings: number;
	
	    static createFrom(source: any = {}) {
	        return new ServiceReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.top_services = this.convertValues(source["top_services"], NameValuePair);
	        this.least_used = this.convertValues(source["least_used"], NameValuePair);
	        this.revenue_by_service = this.convertValues(source["revenue_by_service"], NameValuePair);
	        this.avg_service_value = source["avg_service_value"];
	        this.total_bookings = source["total_bookings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Staff {
	    id: number[];
	    staff_code: string;
	    full_name: string;
	    phone: string;
	    email: string;
	    gender: string;
	    designation: string;
	    // Go type: time
	    joining_date: any;
	    base_salary: number;
	    commission_percentage: number;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Staff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.staff_code = source["staff_code"];
	        this.full_name = source["full_name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gender = source["gender"];
	        this.designation = source["designation"];
	        this.joining_date = this.convertValues(source["joining_date"], null);
	        this.base_salary = source["base_salary"];
	        this.commission_percentage = source["commission_percentage"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StaffPerformanceDaily {
	    id: number[];
	    staff_id: number[];
	    business_date: string;
	    invoice_count: number;
	    customer_count: number;
	    service_count: number;
	    revenue: number;
	    commission_amount: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new StaffPerformanceDaily(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.staff_id = source["staff_id"];
	        this.business_date = source["business_date"];
	        this.invoice_count = source["invoice_count"];
	        this.customer_count = source["customer_count"];
	        this.service_count = source["service_count"];
	        this.revenue = source["revenue"];
	        this.commission_amount = source["commission_amount"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StaffPerformanceSummary {
	    staff_id: number[];
	    staff_name: string;
	    revenue: number;
	    customer_count: number;
	    invoice_count: number;
	    service_count: number;
	    avg_bill: number;
	    commission: number;
	    rank: number;
	
	    static createFrom(source: any = {}) {
	        return new StaffPerformanceSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.staff_name = source["staff_name"];
	        this.revenue = source["revenue"];
	        this.customer_count = source["customer_count"];
	        this.invoice_count = source["invoice_count"];
	        this.service_count = source["service_count"];
	        this.avg_bill = source["avg_bill"];
	        this.commission = source["commission"];
	        this.rank = source["rank"];
	    }
	}
	export class StaffReport {
	    top_performers: NameValuePair[];
	    revenue_by_staff: NameValuePair[];
	    customers_by_staff: NameValuePair[];
	    commission_earned: NameValuePair[];
	    salary_cost: number;
	
	    static createFrom(source: any = {}) {
	        return new StaffReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.top_performers = this.convertValues(source["top_performers"], NameValuePair);
	        this.revenue_by_staff = this.convertValues(source["revenue_by_staff"], NameValuePair);
	        this.customers_by_staff = this.convertValues(source["customers_by_staff"], NameValuePair);
	        this.commission_earned = this.convertValues(source["commission_earned"], NameValuePair);
	        this.salary_cost = source["salary_cost"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StockTransaction {
	    id: number[];
	    product_id: number[];
	    product_name?: string;
	    transaction_type: string;
	    quantity: number;
	    unit_cost: number;
	    reference_type: string;
	    reference_id: string;
	    notes: string;
	    transaction_date: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new StockTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.transaction_type = source["transaction_type"];
	        this.quantity = source["quantity"];
	        this.unit_cost = source["unit_cost"];
	        this.reference_type = source["reference_type"];
	        this.reference_id = source["reference_id"];
	        this.notes = source["notes"];
	        this.transaction_date = source["transaction_date"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TaxRate {
	    id: number[];
	    name: string;
	    hsn_code: string;
	    cgst_rate: number;
	    sgst_rate: number;
	    igst_rate: number;
	    category: string;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TaxRate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.hsn_code = source["hsn_code"];
	        this.cgst_rate = source["cgst_rate"];
	        this.sgst_rate = source["sgst_rate"];
	        this.igst_rate = source["igst_rate"];
	        this.category = source["category"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class UpdateRecord {
	    id: number[];
	    from_version: string;
	    to_version: string;
	    update_date: string;
	    status: string;
	    error_message?: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.from_version = source["from_version"];
	        this.to_version = source["to_version"];
	        this.update_date = source["update_date"];
	        this.status = source["status"];
	        this.error_message = source["error_message"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateStatus {
	    current_version: string;
	    latest_version?: string;
	    update_available: boolean;
	    status: string;
	    release_notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current_version = source["current_version"];
	        this.latest_version = source["latest_version"];
	        this.update_available = source["update_available"];
	        this.status = source["status"];
	        this.release_notes = source["release_notes"];
	    }
	}
	export class WAMessageStats {
	    total_sent: number;
	    delivered: number;
	    read: number;
	    failed: number;
	    queued: number;
	
	    static createFrom(source: any = {}) {
	        return new WAMessageStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_sent = source["total_sent"];
	        this.delivered = source["delivered"];
	        this.read = source["read"];
	        this.failed = source["failed"];
	        this.queued = source["queued"];
	    }
	}
	export class WhatsAppMessage {
	    id: number[];
	    template_id: string;
	    recipient_phone: string;
	    recipient_name: string;
	    message_body: string;
	    status: string;
	    provider: string;
	    provider_message_id: string;
	    error_message: string;
	    sent_at: string;
	    delivered_at: string;
	    read_at: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new WhatsAppMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.template_id = source["template_id"];
	        this.recipient_phone = source["recipient_phone"];
	        this.recipient_name = source["recipient_name"];
	        this.message_body = source["message_body"];
	        this.status = source["status"];
	        this.provider = source["provider"];
	        this.provider_message_id = source["provider_message_id"];
	        this.error_message = source["error_message"];
	        this.sent_at = source["sent_at"];
	        this.delivered_at = source["delivered_at"];
	        this.read_at = source["read_at"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WhatsAppTemplate {
	    id: number[];
	    name: string;
	    category: string;
	    body: string;
	    variables: string;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new WhatsAppTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.body = source["body"];
	        this.variables = source["variables"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace ports {
	
	export class CommissionStaffSummary {
	    staff_id: number[];
	    staff_name: string;
	    revenue: number;
	    commission: number;
	
	    static createFrom(source: any = {}) {
	        return new CommissionStaffSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.staff_name = source["staff_name"];
	        this.revenue = source["revenue"];
	        this.commission = source["commission"];
	    }
	}

}

export namespace usecase {
	
	export class CommissionStats {
	    total_commission_this_month: number;
	    top_earner?: ports.CommissionStaffSummary;
	    avg_commission: number;
	
	    static createFrom(source: any = {}) {
	        return new CommissionStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_commission_this_month = source["total_commission_this_month"];
	        this.top_earner = this.convertValues(source["top_earner"], ports.CommissionStaffSummary);
	        this.avg_commission = source["avg_commission"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateAdvanceInput {
	    staff_id: string;
	    amount: number;
	    advance_date: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateAdvanceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.amount = source["amount"];
	        this.advance_date = source["advance_date"];
	        this.reason = source["reason"];
	    }
	}
	export class CreateCustomerInput {
	    full_name: string;
	    phone: string;
	    email: string;
	    gender: string;
	    date_of_birth: string;
	    anniversary_date: string;
	    address: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCustomerInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full_name = source["full_name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gender = source["gender"];
	        this.date_of_birth = source["date_of_birth"];
	        this.anniversary_date = source["anniversary_date"];
	        this.address = source["address"];
	        this.notes = source["notes"];
	    }
	}
	export class CreateExpenseInput {
	    category_id: string;
	    amount: number;
	    expense_date: string;
	    payment_method: string;
	    vendor_name: string;
	    invoice_reference: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateExpenseInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category_id = source["category_id"];
	        this.amount = source["amount"];
	        this.expense_date = source["expense_date"];
	        this.payment_method = source["payment_method"];
	        this.vendor_name = source["vendor_name"];
	        this.invoice_reference = source["invoice_reference"];
	        this.description = source["description"];
	    }
	}
	export class CreateInvoiceItemInput {
	    service_id: string;
	    quantity: number;
	    discount: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateInvoiceItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.service_id = source["service_id"];
	        this.quantity = source["quantity"];
	        this.discount = source["discount"];
	    }
	}
	export class CreateInvoiceInput {
	    customer_id: string;
	    staff_id: string;
	    items: CreateInvoiceItemInput[];
	    discount: number;
	    tax: number;
	    payment_method: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateInvoiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.customer_id = source["customer_id"];
	        this.staff_id = source["staff_id"];
	        this.items = this.convertValues(source["items"], CreateInvoiceItemInput);
	        this.discount = source["discount"];
	        this.tax = source["tax"];
	        this.payment_method = source["payment_method"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class CreateProductInput {
	    name: string;
	    category: string;
	    brand: string;
	    unit: string;
	    sku: string;
	    purchase_price: number;
	    selling_price: number;
	    minimum_stock: number;
	    maximum_stock: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateProductInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.brand = source["brand"];
	        this.unit = source["unit"];
	        this.sku = source["sku"];
	        this.purchase_price = source["purchase_price"];
	        this.selling_price = source["selling_price"];
	        this.minimum_stock = source["minimum_stock"];
	        this.maximum_stock = source["maximum_stock"];
	    }
	}
	export class PurchaseItemInput {
	    product_id: string;
	    quantity: number;
	    unit_price: number;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	    }
	}
	export class CreatePurchaseInput {
	    vendor_name: string;
	    invoice_number: string;
	    purchase_date: string;
	    notes: string;
	    items: PurchaseItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreatePurchaseInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vendor_name = source["vendor_name"];
	        this.invoice_number = source["invoice_number"];
	        this.purchase_date = source["purchase_date"];
	        this.notes = source["notes"];
	        this.items = this.convertValues(source["items"], PurchaseItemInput);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateRuleInput {
	    rule_name: string;
	    rule_type: string;
	    target_type: string;
	    target_id: string;
	    calculation_type: string;
	    calculation_value: number;
	    minimum_target: number;
	    maximum_target: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateRuleInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rule_name = source["rule_name"];
	        this.rule_type = source["rule_type"];
	        this.target_type = source["target_type"];
	        this.target_id = source["target_id"];
	        this.calculation_type = source["calculation_type"];
	        this.calculation_value = source["calculation_value"];
	        this.minimum_target = source["minimum_target"];
	        this.maximum_target = source["maximum_target"];
	    }
	}
	export class CreateServiceInput {
	    name: string;
	    category: string;
	    description: string;
	    duration_minutes: number;
	    price: number;
	    cost_price: number;
	    commission_type: string;
	    commission_value: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateServiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.duration_minutes = source["duration_minutes"];
	        this.price = source["price"];
	        this.cost_price = source["cost_price"];
	        this.commission_type = source["commission_type"];
	        this.commission_value = source["commission_value"];
	    }
	}
	export class CreateStaffInput {
	    full_name: string;
	    phone: string;
	    email: string;
	    gender: string;
	    designation: string;
	    joining_date: string;
	    base_salary: number;
	    commission_percentage: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateStaffInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full_name = source["full_name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gender = source["gender"];
	        this.designation = source["designation"];
	        this.joining_date = source["joining_date"];
	        this.base_salary = source["base_salary"];
	        this.commission_percentage = source["commission_percentage"];
	    }
	}
	export class CustomerStats {
	    total: number;
	    active: number;
	    inactive: number;
	    new_this_month: number;
	    birthday_today: number;
	
	    static createFrom(source: any = {}) {
	        return new CustomerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.active = source["active"];
	        this.inactive = source["inactive"];
	        this.new_this_month = source["new_this_month"];
	        this.birthday_today = source["birthday_today"];
	    }
	}
	export class DailyPerformanceInput {
	    staff_id: string;
	    date: string;
	
	    static createFrom(source: any = {}) {
	        return new DailyPerformanceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.date = source["date"];
	    }
	}
	export class ExpenseReport {
	    date_from: string;
	    date_to: string;
	    total_expenses: number;
	    expenses_by_category: domain.CategoryExpense[];
	    expense_count: number;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date_from = source["date_from"];
	        this.date_to = source["date_to"];
	        this.total_expenses = source["total_expenses"];
	        this.expenses_by_category = this.convertValues(source["expenses_by_category"], domain.CategoryExpense);
	        this.expense_count = source["expense_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExpenseReportInput {
	    DateFrom: string;
	    DateTo: string;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseReportInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DateFrom = source["DateFrom"];
	        this.DateTo = source["DateTo"];
	    }
	}
	export class GenerateSalaryInput {
	    month: number;
	    year: number;
	
	    static createFrom(source: any = {}) {
	        return new GenerateSalaryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.year = source["year"];
	    }
	}
	export class GenerateSalaryOutput {
	    cycle?: domain.SalaryCycle;
	    records: domain.SalaryRecord[];
	
	    static createFrom(source: any = {}) {
	        return new GenerateSalaryOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cycle = this.convertValues(source["cycle"], domain.SalaryCycle);
	        this.records = this.convertValues(source["records"], domain.SalaryRecord);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetStaffCommissionInput {
	    staff_id: string;
	    date_from: string;
	    date_to: string;
	
	    static createFrom(source: any = {}) {
	        return new GetStaffCommissionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.date_from = source["date_from"];
	        this.date_to = source["date_to"];
	    }
	}
	export class InvoiceStats {
	    today_revenue: number;
	    today_invoices: number;
	    avg_bill_value: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.today_revenue = source["today_revenue"];
	        this.today_invoices = source["today_invoices"];
	        this.avg_bill_value = source["avg_bill_value"];
	    }
	}
	export class ListAdvancesInput {
	    staff_id: string;
	    status: string;
	    page: number;
	    per_page: number;
	
	    static createFrom(source: any = {}) {
	        return new ListAdvancesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.status = source["status"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	    }
	}
	export class ListAdvancesOutput {
	    advances: domain.Advance[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListAdvancesOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.advances = this.convertValues(source["advances"], domain.Advance);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListCustomerInput {
	    search: string;
	    status: string;
	    page: number;
	    per_page: number;
	
	    static createFrom(source: any = {}) {
	        return new ListCustomerInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search = source["search"];
	        this.status = source["status"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	    }
	}
	export class ListCustomerOutput {
	    customers: domain.Customer[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListCustomerOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.customers = this.convertValues(source["customers"], domain.Customer);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListExpensesInput {
	    CategoryID: string;
	    Status: string;
	    PaymentMethod: string;
	    DateFrom: string;
	    DateTo: string;
	    Search: string;
	    Page: number;
	    PerPage: number;
	
	    static createFrom(source: any = {}) {
	        return new ListExpensesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CategoryID = source["CategoryID"];
	        this.Status = source["Status"];
	        this.PaymentMethod = source["PaymentMethod"];
	        this.DateFrom = source["DateFrom"];
	        this.DateTo = source["DateTo"];
	        this.Search = source["Search"];
	        this.Page = source["Page"];
	        this.PerPage = source["PerPage"];
	    }
	}
	export class ListExpensesOutput {
	    expenses: domain.Expense[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListExpensesOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.expenses = this.convertValues(source["expenses"], domain.Expense);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListInvoiceInput {
	    customer_id: string;
	    staff_id: string;
	    payment_status: string;
	    date_from: string;
	    date_to: string;
	    search: string;
	    page: number;
	    per_page: number;
	
	    static createFrom(source: any = {}) {
	        return new ListInvoiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.customer_id = source["customer_id"];
	        this.staff_id = source["staff_id"];
	        this.payment_status = source["payment_status"];
	        this.date_from = source["date_from"];
	        this.date_to = source["date_to"];
	        this.search = source["search"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	    }
	}
	export class ListInvoiceOutput {
	    invoices: domain.Invoice[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListInvoiceOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.invoices = this.convertValues(source["invoices"], domain.Invoice);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListProductsInput {
	    Category: string;
	    Status: string;
	    Search: string;
	    LowStock: boolean;
	    Page: number;
	    PerPage: number;
	
	    static createFrom(source: any = {}) {
	        return new ListProductsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Category = source["Category"];
	        this.Status = source["Status"];
	        this.Search = source["Search"];
	        this.LowStock = source["LowStock"];
	        this.Page = source["Page"];
	        this.PerPage = source["PerPage"];
	    }
	}
	export class ListProductsOutput {
	    products: domain.Product[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListProductsOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.products = this.convertValues(source["products"], domain.Product);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListPurchasesInput {
	    DateFrom: string;
	    DateTo: string;
	    Page: number;
	    PerPage: number;
	
	    static createFrom(source: any = {}) {
	        return new ListPurchasesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DateFrom = source["DateFrom"];
	        this.DateTo = source["DateTo"];
	        this.Page = source["Page"];
	        this.PerPage = source["PerPage"];
	    }
	}
	export class ListPurchasesOutput {
	    purchases: domain.PurchaseEntry[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListPurchasesOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.purchases = this.convertValues(source["purchases"], domain.PurchaseEntry);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListRulesInput {
	    rule_type: string;
	    target_type: string;
	    is_active?: boolean;
	    page: number;
	    per_page: number;
	
	    static createFrom(source: any = {}) {
	        return new ListRulesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rule_type = source["rule_type"];
	        this.target_type = source["target_type"];
	        this.is_active = source["is_active"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	    }
	}
	export class ListRulesOutput {
	    rules: domain.CommissionRule[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListRulesOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rules = this.convertValues(source["rules"], domain.CommissionRule);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListSalariesInput {
	    month: number;
	    year: number;
	
	    static createFrom(source: any = {}) {
	        return new ListSalariesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.year = source["year"];
	    }
	}
	export class ListServiceInput {
	    search: string;
	    status: string;
	    category: string;
	    page: number;
	    per_page: number;
	
	    static createFrom(source: any = {}) {
	        return new ListServiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search = source["search"];
	        this.status = source["status"];
	        this.category = source["category"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	    }
	}
	export class ListServiceOutput {
	    services: domain.Service[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListServiceOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.services = this.convertValues(source["services"], domain.Service);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListStaffInput {
	    search: string;
	    status: string;
	    designation: string;
	    page: number;
	    per_page: number;
	
	    static createFrom(source: any = {}) {
	        return new ListStaffInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search = source["search"];
	        this.status = source["status"];
	        this.designation = source["designation"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	    }
	}
	export class ListStaffOutput {
	    staff: domain.Staff[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListStaffOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff = this.convertValues(source["staff"], domain.Staff);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListStockHistoryInput {
	    ProductID: string;
	    TransactionType: string;
	    DateFrom: string;
	    DateTo: string;
	    Page: number;
	    PerPage: number;
	
	    static createFrom(source: any = {}) {
	        return new ListStockHistoryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ProductID = source["ProductID"];
	        this.TransactionType = source["TransactionType"];
	        this.DateFrom = source["DateFrom"];
	        this.DateTo = source["DateTo"];
	        this.Page = source["Page"];
	        this.PerPage = source["PerPage"];
	    }
	}
	export class ListStockHistoryOutput {
	    transactions: domain.StockTransaction[];
	    total: number;
	    page: number;
	    per_page: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new ListStockHistoryOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.transactions = this.convertValues(source["transactions"], domain.StockTransaction);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.per_page = source["per_page"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MonthlyCommissionInput {
	    month: string;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyCommissionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	    }
	}
	export class PerformanceStats {
	    top_performer_today?: domain.StaffPerformanceSummary;
	    top_performer_month?: domain.StaffPerformanceSummary;
	    total_revenue_today: number;
	    total_customers_today: number;
	    avg_bill_today: number;
	
	    static createFrom(source: any = {}) {
	        return new PerformanceStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.top_performer_today = this.convertValues(source["top_performer_today"], domain.StaffPerformanceSummary);
	        this.top_performer_month = this.convertValues(source["top_performer_month"], domain.StaffPerformanceSummary);
	        this.total_revenue_today = source["total_revenue_today"];
	        this.total_customers_today = source["total_customers_today"];
	        this.avg_bill_today = source["avg_bill_today"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PeriodPerformanceInput {
	    staff_id: string;
	    date_from: string;
	    date_to: string;
	
	    static createFrom(source: any = {}) {
	        return new PeriodPerformanceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.date_from = source["date_from"];
	        this.date_to = source["date_to"];
	    }
	}
	export class ProfitLossInput {
	    DateFrom: string;
	    DateTo: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfitLossInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DateFrom = source["DateFrom"];
	        this.DateTo = source["DateTo"];
	    }
	}
	
	export class RecordPaymentInput {
	    amount: number;
	    payment_method: string;
	    reference_number: string;
	
	    static createFrom(source: any = {}) {
	        return new RecordPaymentInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.amount = source["amount"];
	        this.payment_method = source["payment_method"];
	        this.reference_number = source["reference_number"];
	    }
	}
	export class SalaryStats {
	    total_payroll: number;
	    pending_payments: number;
	    paid_salaries: number;
	    outstanding_advances: number;
	
	    static createFrom(source: any = {}) {
	        return new SalaryStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_payroll = source["total_payroll"];
	        this.pending_payments = source["pending_payments"];
	        this.paid_salaries = source["paid_salaries"];
	        this.outstanding_advances = source["outstanding_advances"];
	    }
	}
	export class ServiceStats {
	    total: number;
	    active: number;
	    inactive: number;
	    avg_price: number;
	
	    static createFrom(source: any = {}) {
	        return new ServiceStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.active = source["active"];
	        this.inactive = source["inactive"];
	        this.avg_price = source["avg_price"];
	    }
	}
	export class StaffCommissionOutput {
	    staff_id: number[];
	    total_revenue: number;
	    commission: number;
	    transactions: domain.CommissionTransaction[];
	
	    static createFrom(source: any = {}) {
	        return new StaffCommissionOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.staff_id = source["staff_id"];
	        this.total_revenue = source["total_revenue"];
	        this.commission = source["commission"];
	        this.transactions = this.convertValues(source["transactions"], domain.CommissionTransaction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StaffStats {
	    total: number;
	    active: number;
	    inactive: number;
	
	    static createFrom(source: any = {}) {
	        return new StaffStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.active = source["active"];
	        this.inactive = source["inactive"];
	    }
	}
	export class StockAdjustInput {
	    product_id: string;
	    transaction_type: string;
	    quantity: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new StockAdjustInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.transaction_type = source["transaction_type"];
	        this.quantity = source["quantity"];
	        this.notes = source["notes"];
	    }
	}
	export class TopPerformersInput {
	    date_from: string;
	    date_to: string;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new TopPerformersInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date_from = source["date_from"];
	        this.date_to = source["date_to"];
	        this.limit = source["limit"];
	    }
	}
	export class UpdateCustomerInput {
	    full_name: string;
	    phone: string;
	    email: string;
	    gender: string;
	    date_of_birth: string;
	    anniversary_date: string;
	    address: string;
	    notes: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCustomerInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full_name = source["full_name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gender = source["gender"];
	        this.date_of_birth = source["date_of_birth"];
	        this.anniversary_date = source["anniversary_date"];
	        this.address = source["address"];
	        this.notes = source["notes"];
	        this.status = source["status"];
	    }
	}
	export class UpdateExpenseInput {
	    category_id: string;
	    amount: number;
	    expense_date: string;
	    payment_method: string;
	    vendor_name: string;
	    invoice_reference: string;
	    description: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateExpenseInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category_id = source["category_id"];
	        this.amount = source["amount"];
	        this.expense_date = source["expense_date"];
	        this.payment_method = source["payment_method"];
	        this.vendor_name = source["vendor_name"];
	        this.invoice_reference = source["invoice_reference"];
	        this.description = source["description"];
	        this.status = source["status"];
	    }
	}
	export class UpdateProductInput {
	    name: string;
	    category: string;
	    brand: string;
	    unit: string;
	    sku: string;
	    purchase_price: number;
	    selling_price: number;
	    minimum_stock: number;
	    maximum_stock: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProductInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.brand = source["brand"];
	        this.unit = source["unit"];
	        this.sku = source["sku"];
	        this.purchase_price = source["purchase_price"];
	        this.selling_price = source["selling_price"];
	        this.minimum_stock = source["minimum_stock"];
	        this.maximum_stock = source["maximum_stock"];
	        this.status = source["status"];
	    }
	}
	export class UpdateRuleInput {
	    rule_name: string;
	    rule_type: string;
	    target_type: string;
	    target_id: string;
	    calculation_type: string;
	    calculation_value: number;
	    minimum_target: number;
	    maximum_target: number;
	    is_active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRuleInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rule_name = source["rule_name"];
	        this.rule_type = source["rule_type"];
	        this.target_type = source["target_type"];
	        this.target_id = source["target_id"];
	        this.calculation_type = source["calculation_type"];
	        this.calculation_value = source["calculation_value"];
	        this.minimum_target = source["minimum_target"];
	        this.maximum_target = source["maximum_target"];
	        this.is_active = source["is_active"];
	    }
	}
	export class UpdateServiceInput {
	    name: string;
	    category: string;
	    description: string;
	    duration_minutes: number;
	    price: number;
	    cost_price: number;
	    commission_type: string;
	    commission_value: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateServiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.duration_minutes = source["duration_minutes"];
	        this.price = source["price"];
	        this.cost_price = source["cost_price"];
	        this.commission_type = source["commission_type"];
	        this.commission_value = source["commission_value"];
	        this.status = source["status"];
	    }
	}
	export class UpdateStaffInput {
	    full_name: string;
	    phone: string;
	    email: string;
	    gender: string;
	    designation: string;
	    joining_date: string;
	    base_salary: number;
	    commission_percentage: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateStaffInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full_name = source["full_name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gender = source["gender"];
	        this.designation = source["designation"];
	        this.joining_date = source["joining_date"];
	        this.base_salary = source["base_salary"];
	        this.commission_percentage = source["commission_percentage"];
	        this.status = source["status"];
	    }
	}

}

