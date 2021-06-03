package servicenow

type CreateCMDBInstance_Response struct {
	Result struct {
		OutboundRelations []struct {
			SysID string `json:"sys_id"`
			Type  struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"type"`
			Target struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"target"`
		} `json:"outbound_relations"`
		Attributes struct {
			FirewallStatus       string `json:"firewall_status"`
			OsAddressWidth       string `json:"os_address_width"`
			AttestedDate         string `json:"attested_date"`
			OperationalStatus    string `json:"operational_status"`
			OsServicePack        string `json:"os_service_pack"`
			CPUCoreThread        string `json:"cpu_core_thread"`
			CPUManufacturer      string `json:"cpu_manufacturer"`
			SysUpdatedOn         string `json:"sys_updated_on"`
			DiscoverySource      string `json:"discovery_source"`
			FirstDiscovered      string `json:"first_discovered"`
			DueIn                string `json:"due_in"`
			UsedFor              string `json:"used_for"`
			InvoiceNumber        string `json:"invoice_number"`
			GlAccount            string `json:"gl_account"`
			SysCreatedBy         string `json:"sys_created_by"`
			WarrantyExpiration   string `json:"warranty_expiration"`
			RAM                  string `json:"ram"`
			CPUName              string `json:"cpu_name"`
			CPUSpeed             string `json:"cpu_speed"`
			OwnedBy              string `json:"owned_by"`
			CheckedOut           string `json:"checked_out"`
			KernelRelease        string `json:"kernel_release"`
			SysDomainPath        string `json:"sys_domain_path"`
			Classification       string `json:"classification"`
			DiskSpace            string `json:"disk_space"`
			ObjectID             string `json:"object_id"`
			MaintenanceSchedule  string `json:"maintenance_schedule"`
			CostCenter           string `json:"cost_center"`
			AttestedBy           string `json:"attested_by"`
			DNSDomain            string `json:"dns_domain"`
			Assigned             string `json:"assigned"`
			PurchaseDate         string `json:"purchase_date"`
			LifeCycleStage       string `json:"life_cycle_stage"`
			ShortDescription     string `json:"short_description"`
			CdSpeed              string `json:"cd_speed"`
			Floppy               string `json:"floppy"`
			ManagedBy            string `json:"managed_by"`
			OsDomain             string `json:"os_domain"`
			LastDiscovered       string `json:"last_discovered"`
			CanPrint             string `json:"can_print"`
			SysClassName         string `json:"sys_class_name"`
			Manufacturer         string `json:"manufacturer"`
			CPUCount             string `json:"cpu_count"`
			Vendor               string `json:"vendor"`
			LifeCycleStageStatus string `json:"life_cycle_stage_status"`
			ModelNumber          string `json:"model_number"`
			AssignedTo           string `json:"assigned_to"`
			StartDate            string `json:"start_date"`
			OsVersion            string `json:"os_version"`
			SerialNumber         string `json:"serial_number"`
			CdRom                string `json:"cd_rom"`
			SupportGroup         string `json:"support_group"`
			Unverified           string `json:"unverified"`
			CorrelationID        string `json:"correlation_id"`
			Attributes           string `json:"attributes"`
			Asset                string `json:"asset"`
			FormFactor           string `json:"form_factor"`
			CPUCoreCount         string `json:"cpu_core_count"`
			SkipSync             string `json:"skip_sync"`
			AttestationScore     string `json:"attestation_score"`
			SysUpdatedBy         string `json:"sys_updated_by"`
			SysCreatedOn         string `json:"sys_created_on"`
			SysDomain            struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"sys_domain"`
			CPUType           string `json:"cpu_type"`
			InstallDate       string `json:"install_date"`
			AssetTag          string `json:"asset_tag"`
			DrBackup          string `json:"dr_backup"`
			HardwareSubstatus string `json:"hardware_substatus"`
			Fqdn              string `json:"fqdn"`
			ChangeControl     string `json:"change_control"`
			InternetFacing    string `json:"internet_facing"`
			DeliveryDate      string `json:"delivery_date"`
			HardwareStatus    string `json:"hardware_status"`
			InstallStatus     string `json:"install_status"`
			SupportedBy       string `json:"supported_by"`
			Name              string `json:"name"`
			Subcategory       string `json:"subcategory"`
			DefaultGateway    string `json:"default_gateway"`
			ChassisType       string `json:"chassis_type"`
			Virtual           string `json:"virtual"`
			AssignmentGroup   string `json:"assignment_group"`
			ManagedByGroup    string `json:"managed_by_group"`
			SysID             string `json:"sys_id"`
			PoNumber          string `json:"po_number"`
			CheckedIn         string `json:"checked_in"`
			SysClassPath      string `json:"sys_class_path"`
			MacAddress        string `json:"mac_address"`
			Company           string `json:"company"`
			Justification     string `json:"justification"`
			Department        string `json:"department"`
			Cost              string `json:"cost"`
			Comments          string `json:"comments"`
			Os                string `json:"os"`
			SysModCount       string `json:"sys_mod_count"`
			Monitor           string `json:"monitor"`
			ModelID           string `json:"model_id"`
			IPAddress         string `json:"ip_address"`
			DuplicateOf       string `json:"duplicate_of"`
			SysTags           string `json:"sys_tags"`
			CostCc            string `json:"cost_cc"`
			OrderDate         string `json:"order_date"`
			Schedule          string `json:"schedule"`
			Environment       string `json:"environment"`
			Due               string `json:"due"`
			Attested          string `json:"attested"`
			Location          string `json:"location"`
			Category          string `json:"category"`
			FaultCount        string `json:"fault_count"`
			HostName          string `json:"host_name"`
			LeaseID           string `json:"lease_id"`
		} `json:"attributes"`
		InboundRelations []struct {
			SysID string `json:"sys_id"`
			Type  struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"type"`
			Target struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"target"`
		} `json:"inbound_relations"`
	} `json:"result"`
}

type UpdateCMDBInstance_Response struct {
	Result struct {
		OutboundRelations []struct {
			SysID string `json:"sys_id"`
			Type  struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"type"`
			Target struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"target"`
		} `json:"outbound_relations"`
		Attributes struct {
			FirewallStatus       string `json:"firewall_status"`
			OsAddressWidth       string `json:"os_address_width"`
			AttestedDate         string `json:"attested_date"`
			OperationalStatus    string `json:"operational_status"`
			OsServicePack        string `json:"os_service_pack"`
			CPUCoreThread        string `json:"cpu_core_thread"`
			CPUManufacturer      string `json:"cpu_manufacturer"`
			SysUpdatedOn         string `json:"sys_updated_on"`
			DiscoverySource      string `json:"discovery_source"`
			FirstDiscovered      string `json:"first_discovered"`
			DueIn                string `json:"due_in"`
			UsedFor              string `json:"used_for"`
			InvoiceNumber        string `json:"invoice_number"`
			GlAccount            string `json:"gl_account"`
			SysCreatedBy         string `json:"sys_created_by"`
			WarrantyExpiration   string `json:"warranty_expiration"`
			RAM                  string `json:"ram"`
			CPUName              string `json:"cpu_name"`
			CPUSpeed             string `json:"cpu_speed"`
			OwnedBy              string `json:"owned_by"`
			CheckedOut           string `json:"checked_out"`
			KernelRelease        string `json:"kernel_release"`
			SysDomainPath        string `json:"sys_domain_path"`
			Classification       string `json:"classification"`
			DiskSpace            string `json:"disk_space"`
			ObjectID             string `json:"object_id"`
			MaintenanceSchedule  string `json:"maintenance_schedule"`
			CostCenter           string `json:"cost_center"`
			AttestedBy           string `json:"attested_by"`
			DNSDomain            string `json:"dns_domain"`
			Assigned             string `json:"assigned"`
			PurchaseDate         string `json:"purchase_date"`
			LifeCycleStage       string `json:"life_cycle_stage"`
			ShortDescription     string `json:"short_description"`
			CdSpeed              string `json:"cd_speed"`
			Floppy               string `json:"floppy"`
			ManagedBy            string `json:"managed_by"`
			OsDomain             string `json:"os_domain"`
			LastDiscovered       string `json:"last_discovered"`
			CanPrint             string `json:"can_print"`
			SysClassName         string `json:"sys_class_name"`
			Manufacturer         string `json:"manufacturer"`
			CPUCount             string `json:"cpu_count"`
			Vendor               string `json:"vendor"`
			LifeCycleStageStatus string `json:"life_cycle_stage_status"`
			ModelNumber          string `json:"model_number"`
			AssignedTo           string `json:"assigned_to"`
			StartDate            string `json:"start_date"`
			OsVersion            string `json:"os_version"`
			SerialNumber         string `json:"serial_number"`
			CdRom                string `json:"cd_rom"`
			SupportGroup         string `json:"support_group"`
			Unverified           string `json:"unverified"`
			CorrelationID        string `json:"correlation_id"`
			Attributes           string `json:"attributes"`
			Asset                struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"asset"`
			FormFactor       string `json:"form_factor"`
			CPUCoreCount     string `json:"cpu_core_count"`
			SkipSync         string `json:"skip_sync"`
			AttestationScore string `json:"attestation_score"`
			SysUpdatedBy     string `json:"sys_updated_by"`
			SysCreatedOn     string `json:"sys_created_on"`
			SysDomain        struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"sys_domain"`
			CPUType           string `json:"cpu_type"`
			InstallDate       string `json:"install_date"`
			AssetTag          string `json:"asset_tag"`
			DrBackup          string `json:"dr_backup"`
			HardwareSubstatus string `json:"hardware_substatus"`
			Fqdn              string `json:"fqdn"`
			ChangeControl     string `json:"change_control"`
			InternetFacing    string `json:"internet_facing"`
			DeliveryDate      string `json:"delivery_date"`
			HardwareStatus    string `json:"hardware_status"`
			InstallStatus     string `json:"install_status"`
			SupportedBy       string `json:"supported_by"`
			Name              string `json:"name"`
			Subcategory       string `json:"subcategory"`
			DefaultGateway    string `json:"default_gateway"`
			ChassisType       string `json:"chassis_type"`
			Virtual           string `json:"virtual"`
			AssignmentGroup   string `json:"assignment_group"`
			ManagedByGroup    string `json:"managed_by_group"`
			SysID             string `json:"sys_id"`
			PoNumber          string `json:"po_number"`
			CheckedIn         string `json:"checked_in"`
			SysClassPath      string `json:"sys_class_path"`
			MacAddress        string `json:"mac_address"`
			Company           string `json:"company"`
			Justification     string `json:"justification"`
			Department        string `json:"department"`
			Cost              string `json:"cost"`
			Comments          string `json:"comments"`
			Os                string `json:"os"`
			SysModCount       string `json:"sys_mod_count"`
			Monitor           string `json:"monitor"`
			ModelID           struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"model_id"`
			IPAddress   string `json:"ip_address"`
			DuplicateOf string `json:"duplicate_of"`
			SysTags     string `json:"sys_tags"`
			CostCc      string `json:"cost_cc"`
			OrderDate   string `json:"order_date"`
			Schedule    string `json:"schedule"`
			Environment string `json:"environment"`
			Due         string `json:"due"`
			Attested    string `json:"attested"`
			Location    string `json:"location"`
			Category    string `json:"category"`
			FaultCount  string `json:"fault_count"`
			HostName    string `json:"host_name"`
			LeaseID     string `json:"lease_id"`
		} `json:"attributes"`
		InboundRelations []struct {
			SysID string `json:"sys_id"`
			Type  struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"type"`
			Target struct {
				DisplayValue string `json:"display_value"`
				Link         string `json:"link"`
				Value        string `json:"value"`
			} `json:"target"`
		} `json:"inbound_relations"`
	} `json:"result"`
}

type CreateIncident_Response struct {
	Result struct {
		UponApproval           string `json:"upon_approval"`
		Location               string `json:"location"`
		ExpectedStart          string `json:"expected_start"`
		ReopenCount            string `json:"reopen_count"`
		CloseNotes             string `json:"close_notes"`
		AdditionalAssigneeList string `json:"additional_assignee_list"`
		Impact                 string `json:"impact"`
		Urgency                string `json:"urgency"`
		CorrelationID          string `json:"correlation_id"`
		SysTags                string `json:"sys_tags"`
		SysDomain              struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"sys_domain"`
		Description        string `json:"description"`
		GroupList          string `json:"group_list"`
		Priority           string `json:"priority"`
		DeliveryPlan       string `json:"delivery_plan"`
		SysModCount        string `json:"sys_mod_count"`
		WorkNotesList      string `json:"work_notes_list"`
		BusinessService    string `json:"business_service"`
		FollowUp           string `json:"follow_up"`
		ClosedAt           string `json:"closed_at"`
		SLADue             string `json:"sla_due"`
		DeliveryTask       string `json:"delivery_task"`
		SysUpdatedOn       string `json:"sys_updated_on"`
		Parent             string `json:"parent"`
		WorkEnd            string `json:"work_end"`
		Number             string `json:"number"`
		ClosedBy           string `json:"closed_by"`
		WorkStart          string `json:"work_start"`
		CalendarStc        string `json:"calendar_stc"`
		Category           string `json:"category"`
		BusinessDuration   string `json:"business_duration"`
		IncidentState      string `json:"incident_state"`
		ActivityDue        string `json:"activity_due"`
		CorrelationDisplay string `json:"correlation_display"`
		Company            string `json:"company"`
		Active             string `json:"active"`
		DueDate            string `json:"due_date"`
		AssignmentGroup    struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"assignment_group"`
		CallerID             string `json:"caller_id"`
		Knowledge            string `json:"knowledge"`
		MadeSLA              string `json:"made_sla"`
		CommentsAndWorkNotes string `json:"comments_and_work_notes"`
		ParentIncident       string `json:"parent_incident"`
		State                string `json:"state"`
		UserInput            string `json:"user_input"`
		SysCreatedOn         string `json:"sys_created_on"`
		ApprovalSet          string `json:"approval_set"`
		ReassignmentCount    string `json:"reassignment_count"`
		Rfc                  string `json:"rfc"`
		ChildIncidents       string `json:"child_incidents"`
		OpenedAt             string `json:"opened_at"`
		ShortDescription     string `json:"short_description"`
		Order                string `json:"order"`
		SysUpdatedBy         string `json:"sys_updated_by"`
		ResolvedBy           string `json:"resolved_by"`
		Notify               string `json:"notify"`
		UponReject           string `json:"upon_reject"`
		ApprovalHistory      string `json:"approval_history"`
		ProblemID            string `json:"problem_id"`
		WorkNotes            string `json:"work_notes"`
		CalendarDuration     string `json:"calendar_duration"`
		CloseCode            string `json:"close_code"`
		SysID                string `json:"sys_id"`
		Approval             string `json:"approval"`
		CausedBy             string `json:"caused_by"`
		Severity             string `json:"severity"`
		SysCreatedBy         string `json:"sys_created_by"`
		ResolvedAt           string `json:"resolved_at"`
		AssignedTo           string `json:"assigned_to"`
		BusinessStc          string `json:"business_stc"`
		WfActivity           string `json:"wf_activity"`
		SysDomainPath        string `json:"sys_domain_path"`
		CmdbCi               string `json:"cmdb_ci"`
		OpenedBy             struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"opened_by"`
		Subcategory   string `json:"subcategory"`
		RejectionGoto string `json:"rejection_goto"`
		SysClassName  string `json:"sys_class_name"`
		WatchList     string `json:"watch_list"`
		TimeWorked    string `json:"time_worked"`
		ContactType   string `json:"contact_type"`
		Escalation    string `json:"escalation"`
		Comments      string `json:"comments"`
	} `json:"result"`
}

type UpdateIncident_Response struct {
	Result struct {
		UponApproval string `json:"upon_approval"`
		Location     struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"location"`
		ExpectedStart          string `json:"expected_start"`
		ReopenCount            string `json:"reopen_count"`
		CloseNotes             string `json:"close_notes"`
		AdditionalAssigneeList string `json:"additional_assignee_list"`
		Impact                 string `json:"impact"`
		Urgency                string `json:"urgency"`
		CorrelationID          string `json:"correlation_id"`
		SysTags                string `json:"sys_tags"`
		SysDomain              struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"sys_domain"`
		Description        string `json:"description"`
		GroupList          string `json:"group_list"`
		Priority           string `json:"priority"`
		DeliveryPlan       string `json:"delivery_plan"`
		SysModCount        string `json:"sys_mod_count"`
		WorkNotesList      string `json:"work_notes_list"`
		BusinessService    string `json:"business_service"`
		FollowUp           string `json:"follow_up"`
		ClosedAt           string `json:"closed_at"`
		SLADue             string `json:"sla_due"`
		DeliveryTask       string `json:"delivery_task"`
		SysUpdatedOn       string `json:"sys_updated_on"`
		Parent             string `json:"parent"`
		WorkEnd            string `json:"work_end"`
		Number             string `json:"number"`
		ClosedBy           string `json:"closed_by"`
		WorkStart          string `json:"work_start"`
		CalendarStc        string `json:"calendar_stc"`
		Category           string `json:"category"`
		BusinessDuration   string `json:"business_duration"`
		IncidentState      string `json:"incident_state"`
		ActivityDue        string `json:"activity_due"`
		CorrelationDisplay string `json:"correlation_display"`
		Company            struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"company"`
		Active          string `json:"active"`
		DueDate         string `json:"due_date"`
		AssignmentGroup struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"assignment_group"`
		CallerID struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"caller_id"`
		Knowledge            string `json:"knowledge"`
		MadeSLA              string `json:"made_sla"`
		CommentsAndWorkNotes string `json:"comments_and_work_notes"`
		ParentIncident       string `json:"parent_incident"`
		State                string `json:"state"`
		UserInput            string `json:"user_input"`
		SysCreatedOn         string `json:"sys_created_on"`
		ApprovalSet          string `json:"approval_set"`
		ReassignmentCount    string `json:"reassignment_count"`
		Rfc                  string `json:"rfc"`
		ChildIncidents       string `json:"child_incidents"`
		OpenedAt             string `json:"opened_at"`
		ShortDescription     string `json:"short_description"`
		Order                string `json:"order"`
		SysUpdatedBy         string `json:"sys_updated_by"`
		ResolvedBy           string `json:"resolved_by"`
		Notify               string `json:"notify"`
		UponReject           string `json:"upon_reject"`
		ApprovalHistory      string `json:"approval_history"`
		ProblemID            string `json:"problem_id"`
		WorkNotes            string `json:"work_notes"`
		CalendarDuration     string `json:"calendar_duration"`
		CloseCode            string `json:"close_code"`
		SysID                string `json:"sys_id"`
		Approval             string `json:"approval"`
		CausedBy             string `json:"caused_by"`
		Severity             string `json:"severity"`
		SysCreatedBy         string `json:"sys_created_by"`
		ResolvedAt           string `json:"resolved_at"`
		AssignedTo           struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"assigned_to"`
		BusinessStc   string `json:"business_stc"`
		WfActivity    string `json:"wf_activity"`
		SysDomainPath string `json:"sys_domain_path"`
		CmdbCi        struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"cmdb_ci"`
		OpenedBy struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"opened_by"`
		Subcategory   string `json:"subcategory"`
		RejectionGoto string `json:"rejection_goto"`
		SysClassName  string `json:"sys_class_name"`
		WatchList     string `json:"watch_list"`
		TimeWorked    string `json:"time_worked"`
		ContactType   string `json:"contact_type"`
		Escalation    string `json:"escalation"`
		Comments      string `json:"comments"`
	} `json:"result"`
}

type GetIncident_Response struct {
	Result struct {
		UponApproval string `json:"upon_approval"`
		Location     struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"location"`
		ExpectedStart          string `json:"expected_start"`
		ReopenCount            string `json:"reopen_count"`
		CloseNotes             string `json:"close_notes"`
		AdditionalAssigneeList string `json:"additional_assignee_list"`
		Impact                 string `json:"impact"`
		Urgency                string `json:"urgency"`
		CorrelationID          string `json:"correlation_id"`
		SysTags                string `json:"sys_tags"`
		SysDomain              struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"sys_domain"`
		Description        string `json:"description"`
		GroupList          string `json:"group_list"`
		Priority           string `json:"priority"`
		DeliveryPlan       string `json:"delivery_plan"`
		SysModCount        string `json:"sys_mod_count"`
		WorkNotesList      string `json:"work_notes_list"`
		BusinessService    string `json:"business_service"`
		FollowUp           string `json:"follow_up"`
		ClosedAt           string `json:"closed_at"`
		SLADue             string `json:"sla_due"`
		DeliveryTask       string `json:"delivery_task"`
		SysUpdatedOn       string `json:"sys_updated_on"`
		Parent             string `json:"parent"`
		WorkEnd            string `json:"work_end"`
		Number             string `json:"number"`
		ClosedBy           string `json:"closed_by"`
		WorkStart          string `json:"work_start"`
		CalendarStc        string `json:"calendar_stc"`
		Category           string `json:"category"`
		BusinessDuration   string `json:"business_duration"`
		IncidentState      string `json:"incident_state"`
		ActivityDue        string `json:"activity_due"`
		CorrelationDisplay string `json:"correlation_display"`
		Company            string `json:"company"`
		Active             string `json:"active"`
		DueDate            string `json:"due_date"`
		AssignmentGroup    struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"assignment_group"`
		CallerID struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"caller_id"`
		Knowledge            string `json:"knowledge"`
		MadeSLA              string `json:"made_sla"`
		CommentsAndWorkNotes string `json:"comments_and_work_notes"`
		ParentIncident       string `json:"parent_incident"`
		State                string `json:"state"`
		UserInput            string `json:"user_input"`
		SysCreatedOn         string `json:"sys_created_on"`
		ApprovalSet          string `json:"approval_set"`
		ReassignmentCount    string `json:"reassignment_count"`
		Rfc                  string `json:"rfc"`
		ChildIncidents       string `json:"child_incidents"`
		OpenedAt             string `json:"opened_at"`
		ShortDescription     string `json:"short_description"`
		Order                string `json:"order"`
		SysUpdatedBy         string `json:"sys_updated_by"`
		ResolvedBy           string `json:"resolved_by"`
		Notify               string `json:"notify"`
		UponReject           string `json:"upon_reject"`
		ApprovalHistory      string `json:"approval_history"`
		ProblemID            struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"problem_id"`
		WorkNotes        string `json:"work_notes"`
		CalendarDuration string `json:"calendar_duration"`
		CloseCode        string `json:"close_code"`
		SysID            string `json:"sys_id"`
		Approval         string `json:"approval"`
		CausedBy         string `json:"caused_by"`
		Severity         string `json:"severity"`
		SysCreatedBy     string `json:"sys_created_by"`
		ResolvedAt       string `json:"resolved_at"`
		AssignedTo       string `json:"assigned_to"`
		BusinessStc      string `json:"business_stc"`
		WfActivity       string `json:"wf_activity"`
		SysDomainPath    string `json:"sys_domain_path"`
		CmdbCi           struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"cmdb_ci"`
		OpenedBy struct {
			Link  string `json:"link"`
			Value string `json:"value"`
		} `json:"opened_by"`
		Subcategory   string `json:"subcategory"`
		RejectionGoto string `json:"rejection_goto"`
		SysClassName  string `json:"sys_class_name"`
		WatchList     string `json:"watch_list"`
		TimeWorked    string `json:"time_worked"`
		ContactType   string `json:"contact_type"`
		Escalation    string `json:"escalation"`
		Comments      string `json:"comments"`
	} `json:"result"`
}

type ChangeRequest_Response struct {
	Result []struct {
		SysID struct {
			Value        string `json:"value"`
			DisplayValue string `json:"display_value"`
		} `json:"sys_id"`
		State struct {
			Value        string `json:"value"`
			DisplayValue string `json:"display_value"`
		} `json:"state"`
	} `json:"result"`
}
