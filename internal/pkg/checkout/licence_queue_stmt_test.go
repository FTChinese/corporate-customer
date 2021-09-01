package checkout

import (
	"github.com/FTChinese/ftacademy/pkg/sq"
	"testing"
)

func TestStmtBulkLicenceQueue(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want sq.BulkInsert
	}{
		{
			name: "SQL for bulk insert licence",
			args: args{
				n: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StmtBulkLicenceQueue(tt.args.n)

			//if  !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("StmtBulkLicenceQueue() = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", got.Build())
		})
	}
}
