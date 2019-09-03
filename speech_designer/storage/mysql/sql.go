package mysql

var (
	SqlQueryNodeFormartByRequest = "SELECT " +
			"s.s_id," +
			"s.s_priority," +
			"n.sn_id," +
			"n.sn_parent," +
			"n.sn_action," +
			"n.sn_argc," +
			"n.sn_argv," +
			"n.sn_file," +
			"n.sn_text," +
			"n.sn_description," +
			"t.snt_trigger " +
		"FROM " +
			"speech AS s " +
			"LEFT JOIN s_node AS n ON s.s_id = n.s_id " +
			"LEFT JOIN sn_trigger AS t ON t.sn_id = n.sn_id " +
		"WHERE " +
			"t.snt_trigger = '%s' " +
			"and n.sn_parent = %d " +
			"and s.s_number = '%s' " +
			"and '%s' BETWEEN s.s_bdate and s.s_edate " +
			"and '%s' BETWEEN s.s_btime and s.s_etime " +
			"ORDER BY t.snt_trigger desc limit 1;"
)
