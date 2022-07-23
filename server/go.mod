module ktrn.com/main

go 1.18

replace ktrn.com/service => ./service

replace ktrn.com/dbhandler => ./dbhandler

require (
	ktrn.com/dbhandler v0.0.0-00010101000000-000000000000
	ktrn.com/service v0.0.0-00010101000000-000000000000
)

require github.com/lib/pq v1.10.6

require github.com/google/uuid v1.3.0 // indirect
