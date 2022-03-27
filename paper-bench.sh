#/bin/bash
go test -benchmem -run=^$ -bench ^BenchmarkAssessVMEvidence$ clouditor.io/clouditor/service/assessment > bench-vm.txt
go test -benchmem -run=^$ -bench ^BenchmarkAssessStorageEvidence$ clouditor.io/clouditor/service/assessment > bench-storage.txt
go test -benchmem -run=^$ -bench ^BenchmarkEvidenceTypes$ clouditor.io/clouditor/service/assessment > bench-types.txt