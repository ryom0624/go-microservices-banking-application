module auth

go 1.16

replace local.packages/hoge => ./../lib/hoge

replace local.packages/errs => ./../lib/errs
replace local.packages/logger => ./../lib/logger


require (
	local.packages/errs v0.0.0-00010101000000-000000000000
	local.packages/hoge v0.0.0-00010101000000-000000000000
	local.packages/logger v0.0.0-00010101000000-000000000000
)
