package compliance.cis_azure.rules.cis_4_5_1

import data.cis_azure.test_data
import data.compliance.policy.azure.data_adapter
import data.lib.test

test_violation {
	eval_fail with input as test_data.generate_azure_asset("azure-document-db-database-account", {"isVirtualNetworkFilterEnabled": false})
	eval_fail with input as test_data.generate_azure_asset("azure-document-db-database-account", {})
}

test_pass {
	eval_pass with input as test_data.generate_azure_asset("azure-document-db-database-account", {"isVirtualNetworkFilterEnabled": true})
}

test_not_evaluated {
	not_eval with input as test_data.not_eval_non_exist_type
}

eval_fail {
	test.assert_fail(finding) with data.benchmark_data_adapter as data_adapter
}

eval_pass {
	test.assert_pass(finding) with data.benchmark_data_adapter as data_adapter
}

not_eval {
	not finding with data.benchmark_data_adapter as data_adapter
}
