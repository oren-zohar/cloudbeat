package compliance.cis_azure.rules.cis_4_3_5

import data.compliance.lib.common
import data.compliance.policy.azure.data_adapter
import future.keywords.if

finding = result if {
	# filter
	data_adapter.is_postgresql_server_db

	# set result
	result := common.generate_result_without_expected(
		common.calculate_result(config_enabled),
		{"Resource": data_adapter.resource},
	)
}

default config_enabled = false

config_enabled if {
	some i
	data_adapter.resource.extension.psqlConfigurations[i].name == "connection_throttling"
	data_adapter.resource.extension.psqlConfigurations[i].value == "on"
}
